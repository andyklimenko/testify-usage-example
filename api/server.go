package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andyklimenko/testify-usage-example/api/entity"
	"github.com/andyklimenko/testify-usage-example/config"
	"github.com/gorilla/mux"
)

type repo interface {
	InsertUser(ctx context.Context, u entity.User) (entity.User, error)
	UserByID(ctx context.Context, id string) (entity.User, error)
	UpdateUser(ctx context.Context, id string, u entity.User) (entity.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userChangelog interface {
	UserCreated(userID string) error
	UserUpdated(userID string) error
	UserDeleted(userID string) error
}

type Server struct {
	httpSrv       *http.Server
	repo          repo
	userChangelog userChangelog
}

func (s *Server) Start() error {
	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-termCh
		if err := s.stop(); err != nil {
			slog.Error("failed to stop the server gracefully", err)
		}
	}()

	if err := s.httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return s.httpSrv.Shutdown(ctx)
}

func setupRouter(s *Server) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users", s.createUser).Methods(http.MethodPost)
	r.HandleFunc("/users/{id}", s.getUser).Methods(http.MethodGet)
	r.HandleFunc("/users/{id}", s.updateUser).Methods(http.MethodPut)
	r.HandleFunc("/users/{id}", s.deleteUser).Methods(http.MethodDelete)

	return r
}

func New(cfg config.Config, s repo, changelog userChangelog) *Server {
	srv := &Server{
		repo:          s,
		userChangelog: changelog,
	}

	srv.httpSrv = &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: setupRouter(srv),
	}

	return srv
}
