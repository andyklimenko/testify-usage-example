package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andyklimenko/testify-usage-example/api/entity"
	"github.com/andyklimenko/testify-usage-example/config"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type repo interface {
	InsertUser(ctx context.Context, u entity.User) (entity.User, error)
	UserByID(ctx context.Context, id string) (entity.User, error)
}

type Server struct {
	httpSrv *http.Server
	repo    repo
}

func (s *Server) Start() error {
	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-termCh
		if err := s.stop(); err != nil {
			logrus.Errorf("failed to stop the server gracefully: %v", err)
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

func New(cfg config.Config, s repo) *Server {
	srv := &Server{
		repo: s,
	}

	srv.httpSrv = &http.Server{
		Addr:    cfg.Server.Addr,
		Handler: setupRouter(srv),
	}

	return srv
}
