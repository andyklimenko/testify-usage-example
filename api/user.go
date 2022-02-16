package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/andyklimenko/testify-usage-example/api/entity"
	"github.com/gorilla/mux"
)

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var u entity.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		s.respondNotOK(w, http.StatusBadRequest, fmt.Errorf("decode request body: %w", err))
	}

	createdUser, err := s.repo.InsertUser(r.Context(), u)
	if err != nil {
		s.respondNotOK(w, http.StatusInternalServerError, fmt.Errorf("create new user: %w", err))
		return
	}

	s.respondOK(w, http.StatusCreated, createdUser)
}

func (s *Server) getUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := mux.Vars(r)["id"]
	if !ok {
		s.respondNotOK(w, http.StatusBadRequest, errors.New("no user id"))
		return
	}

	user, err := s.repo.UserByID(r.Context(), userID)
	if err == nil {
		s.respondOK(w, http.StatusOK, user)
		return
	}

	statusCode := http.StatusInternalServerError
	err = fmt.Errorf("find user by id %s: %w", userID, err)
	if errors.Is(err, entity.ErrNotFound) {
		statusCode = http.StatusNotFound
		err = fmt.Errorf("user %s not found", userID)
	}
	s.respondNotOK(w, statusCode, err)
}

func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	var u entity.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		s.respondNotOK(w, http.StatusBadRequest, fmt.Errorf("decode request body: %w", err))
	}

	userID, ok := mux.Vars(r)["id"]
	if !ok {
		s.respondNotOK(w, http.StatusBadRequest, errors.New("no user id"))
		return
	}

	res, err := s.repo.UpdateUser(r.Context(), userID, u)
	if err == nil {
		s.respondOK(w, http.StatusOK, res)
		return
	}

	statusCode := http.StatusInternalServerError
	err = fmt.Errorf("update user by id %s: %w", userID, err)
	if errors.Is(err, entity.ErrNotFound) {
		statusCode = http.StatusNotFound
		err = fmt.Errorf("user %s not found", userID)
	}
	s.respondNotOK(w, statusCode, err)
}

func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := mux.Vars(r)["id"]
	if !ok {
		s.respondNotOK(w, http.StatusBadRequest, errors.New("no user id"))
		return
	}

	err := s.repo.DeleteUser(r.Context(), userID)
	if err == nil {
		s.respondOK(w, http.StatusOK, nil)
		return
	}

	statusCode := http.StatusInternalServerError
	err = fmt.Errorf("delete user by id %s: %w", userID, err)
	if errors.Is(err, entity.ErrNotFound) {
		statusCode = http.StatusNotFound
		err = fmt.Errorf("user %s not found", userID)
	}
	s.respondNotOK(w, statusCode, err)
}
