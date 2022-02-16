package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"

	"github.com/andyklimenko/testify-usage-example/api/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func closeBody(c io.Closer) {
	if err := c.Close(); err != nil {
		logrus.Errorf("closing response body: %v", err)
	}
}

func (s *srvSuite) setupServer() (string, func()) {
	srv := &Server{
		repo: s.repo,
	}
	testSrv := httptest.NewServer(setupRouter(srv))
	srv.httpSrv = testSrv.Config

	return testSrv.URL, func() {
		testSrv.Close()
	}
}

func (s *srvSuite) TestCreateGetUser() {
	newUser := entity.User{
		FirstName: "John",
		LastName:  "Doe",
	}

	srvURL, closer := s.setupServer()
	defer closer()

	bodyRaw, err := json.Marshal(newUser)
	require.NoError(s.T(), err)

	postUsersResp, err := s.httpCli.Post(srvURL+"/users", "application/json", bytes.NewReader(bodyRaw))
	require.NoError(s.T(), err)

	defer closeBody(postUsersResp.Body)
	require.Equal(s.T(), http.StatusCreated, postUsersResp.StatusCode)

	var userCreated entity.User
	require.NoError(s.T(), json.NewDecoder(postUsersResp.Body).Decode(&userCreated))

	require.Equal(s.T(), "John", userCreated.FirstName)
	require.Equal(s.T(), "Doe", userCreated.LastName)

	// let's get what we've created
	getUsersResp, err := s.httpCli.Get(srvURL + "/users/" + userCreated.ID)
	require.NoError(s.T(), err)

	defer closeBody(getUsersResp.Body)
	require.Equal(s.T(), http.StatusOK, getUsersResp.StatusCode)

	var userGot entity.User
	require.NoError(s.T(), json.NewDecoder(getUsersResp.Body).Decode(&userGot))

	assert.Equal(s.T(), userCreated, userGot)
}

func (s *srvSuite) TestGetMissingUser() {
	srvURL, closer := s.setupServer()
	defer closer()

	missingUserID := uuid.New().String()
	resp, err := s.httpCli.Get(srvURL + "/users/" + missingUserID)
	require.NoError(s.T(), err)

	defer closeBody(resp.Body)
	require.Equal(s.T(), http.StatusNotFound, resp.StatusCode)

	var errResp statusResponse
	require.NoError(s.T(), json.NewDecoder(resp.Body).Decode(&errResp))

	assert.Equal(s.T(), http.StatusNotFound, errResp.Code)
	assert.Equal(s.T(), fmt.Sprintf("user %s not found", missingUserID), errResp.Text)
}
