package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/andyklimenko/testify-usage-example/api/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockedChangelog struct {
	mock.Mock
}

func (m *mockedChangelog) UserCreated(u entity.User) error {
	return m.Called(u).Error(0)
}

func (m *mockedChangelog) UserUpdated(u entity.User) error {
	return m.Called(u).Error(0)
}

func (m *mockedChangelog) UserDeleted(u entity.User) error {
	return m.Called(u).Error(0)
}

func (s *srvSuite) createTestUser(srvURL string, u entity.User) (entity.User, error) {
	bodyRaw, err := json.Marshal(u)
	if err != nil {
		return entity.User{}, err
	}

	postUsersResp, err := s.httpCli.Post(srvURL+"/users", "application/json", bytes.NewReader(bodyRaw))
	if err != nil {
		return entity.User{}, err
	}

	defer entity.CloseBody(postUsersResp.Body)
	if !assert.Equal(s.T(), http.StatusCreated, postUsersResp.StatusCode) {
		return entity.User{}, fmt.Errorf("unexpected status-code: %d", postUsersResp.StatusCode)
	}

	var userCreated entity.User
	if !assert.NoError(s.T(), json.NewDecoder(postUsersResp.Body).Decode(&userCreated)) {
		return entity.User{}, fmt.Errorf("decode response body: %w", err)
	}

	require.Equal(s.T(), u.FirstName, userCreated.FirstName)
	require.Equal(s.T(), u.LastName, userCreated.LastName)

	return userCreated, nil
}

func (s *srvSuite) TestGetUser() {
	var cl mockedChangelog

	newUser := entity.User{
		FirstName: "John",
		LastName:  "Doe",
	}

	srvURL, closer := s.setupServer(&cl)
	defer closer()

	userCh := make(chan entity.User, 1)
	cl.On("UserCreated",
		mock.MatchedBy(func(u entity.User) bool {
			if !assert.Equal(s.T(), "John", u.FirstName) {
				return false
			}

			if !assert.Equal(s.T(), "Doe", u.LastName) {
				return false
			}

			userCh <- u
			return true
		}),
	).Return(nil).Once()

	userCreated, err := s.createTestUser(srvURL, newUser)
	require.NoError(s.T(), err)

	select {
	case <-time.After(time.Second):
		s.T().Fatal("timeout")
	case userFromChangelog := <-userCh:
		require.Equal(s.T(), userCreated.ID, userFromChangelog.ID)
	}

	getUsersResp, err := s.httpCli.Get(srvURL + "/users/" + userCreated.ID)
	require.NoError(s.T(), err)

	defer entity.CloseBody(getUsersResp.Body)
	require.Equal(s.T(), http.StatusOK, getUsersResp.StatusCode)

	var userGot entity.User
	require.NoError(s.T(), json.NewDecoder(getUsersResp.Body).Decode(&userGot))

	assert.Equal(s.T(), userCreated, userGot)

	cl.AssertExpectations(s.T())
}

func (s *srvSuite) TestGetMissingUser() {
	srvURL, closer := s.setupServer(nil)
	defer closer()

	missingUserID := uuid.New().String()
	resp, err := s.httpCli.Get(srvURL + "/users/" + missingUserID)
	require.NoError(s.T(), err)

	defer entity.CloseBody(resp.Body)
	require.Equal(s.T(), http.StatusNotFound, resp.StatusCode)

	var errResp statusResponse
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&errResp))

	assert.Equal(s.T(), http.StatusNotFound, errResp.Code)
	assert.Equal(s.T(), fmt.Sprintf("user %s not found", missingUserID), errResp.Text)
}

func (s *srvSuite) TestUpdateMissingUser() {
	u := entity.User{
		FirstName: "Bob",
		LastName:  "Just Bob",
	}

	srvURL, closer := s.setupServer(nil)
	defer closer()

	missingUserID := uuid.New().String()
	bodyRaw, err := json.Marshal(u)
	require.NoError(s.T(), err)

	req, err := http.NewRequest(http.MethodPut, srvURL+"/users/"+missingUserID, bytes.NewReader(bodyRaw))
	require.NoError(s.T(), err)

	resp, err := s.httpCli.Do(req)
	require.NoError(s.T(), err)

	defer entity.CloseBody(resp.Body)
	require.Equal(s.T(), http.StatusNotFound, resp.StatusCode)

	var errResp statusResponse
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&errResp))

	assert.Equal(s.T(), http.StatusNotFound, errResp.Code)
	assert.Equal(s.T(), fmt.Sprintf("user %s not found", missingUserID), errResp.Text)
}

func (s *srvSuite) TestUpdateUser() {
	var cl mockedChangelog

	newUser := entity.User{
		FirstName: "Anakin",
		LastName:  "Skywalker",
	}

	srvURL, closer := s.setupServer(&cl)
	defer closer()

	cl.On("UserCreated", mock.Anything).Return(nil)
	userCreated, err := s.createTestUser(srvURL, newUser)
	require.NoError(s.T(), err)

	userToUpdate := userCreated
	userToUpdate.FirstName = "Darth"
	userToUpdate.LastName = "Wader"

	bodyRaw, err := json.Marshal(userToUpdate)
	require.NoError(s.T(), err)

	req, err := http.NewRequest(http.MethodPut, srvURL+"/users/"+userCreated.ID, bytes.NewReader(bodyRaw))
	require.NoError(s.T(), err)

	cl.On("UserUpdated",
		mock.MatchedBy(func(u entity.User) bool {
			return assert.Equal(s.T(), "Darth", u.FirstName) &&
				assert.Equal(s.T(), "Wader", u.LastName)
		}),
	).Return(nil).Once()
	getUsersResp, err := s.httpCli.Do(req)
	require.NoError(s.T(), err)

	defer entity.CloseBody(getUsersResp.Body)
	require.Equal(s.T(), http.StatusOK, getUsersResp.StatusCode)

	var userUpdated entity.User
	require.NoError(s.T(), json.NewDecoder(getUsersResp.Body).Decode(&userUpdated))

	assert.Equal(s.T(), userToUpdate, userUpdated)
	cl.AssertExpectations(s.T())
}

func (s *srvSuite) TestDeleteMissingUser() {
	srvURL, closer := s.setupServer(nil)
	defer closer()

	missingUserID := uuid.New().String()
	req, err := http.NewRequest(http.MethodDelete, srvURL+"/users/"+missingUserID, nil)
	require.NoError(s.T(), err)

	resp, err := s.httpCli.Do(req)
	require.NoError(s.T(), err)

	defer entity.CloseBody(resp.Body)
	require.Equal(s.T(), http.StatusNotFound, resp.StatusCode)

	var errResp statusResponse
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&errResp))

	assert.Equal(s.T(), http.StatusNotFound, errResp.Code)
	assert.Equal(s.T(), fmt.Sprintf("user %s not found", missingUserID), errResp.Text)
}

func (s *srvSuite) TestDeleteUser() {
	var cl mockedChangelog

	newUser := entity.User{
		FirstName: "Han",
		LastName:  "Solo",
	}

	srvURL, closer := s.setupServer(&cl)
	defer closer()

	cl.On("UserCreated", mock.Anything).Return(nil).Once()
	userCreated, err := s.createTestUser(srvURL, newUser)
	require.NoError(s.T(), err)

	req, err := http.NewRequest(http.MethodDelete, srvURL+"/users/"+userCreated.ID, nil)
	require.NoError(s.T(), err)

	cl.On("UserDeleted",
		mock.MatchedBy(func(u entity.User) bool {
			return assert.Equal(s.T(), userCreated.ID, u.ID)
		}),
	).Return(nil).Once()
	deleteResp, err := s.httpCli.Do(req)
	require.NoError(s.T(), err)

	defer entity.CloseBody(deleteResp.Body)
	require.Equal(s.T(), http.StatusOK, deleteResp.StatusCode)

	tryToGetOnceAgain, err := http.Get(srvURL + "/users/" + userCreated.ID)
	require.NoError(s.T(), err)
	defer entity.CloseBody(tryToGetOnceAgain.Body)

	var errResp statusResponse
	require.NoError(s.T(), json.NewDecoder(tryToGetOnceAgain.Body).Decode(&errResp))

	// it's really deleted
	assert.Equal(s.T(), http.StatusNotFound, errResp.Code)
	assert.Equal(s.T(), fmt.Sprintf("user %s not found", userCreated.ID), errResp.Text)

	cl.AssertExpectations(s.T())
}
