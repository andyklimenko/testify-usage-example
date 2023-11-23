package api

import (
	"log/slog"

	"github.com/andyklimenko/testify-usage-example/api/entity"
)

func (s *Server) onUserCreated(u entity.User) {
	if err := s.userChangelog.UserCreated(u); err != nil {
		slog.Error("something bad happened while logging user creation", err)
	}
}

func (s *Server) onUserUpdated(u entity.User) {
	if err := s.userChangelog.UserUpdated(u); err != nil {
		slog.Error("something bad happened while logging user update", err)
	}
}

func (s *Server) onUserDeleted(u entity.User) {
	if err := s.userChangelog.UserDeleted(u); err != nil {
		slog.Error("something bad happened while logging user delete", err)
	}
}
