package config

import "errors"

var (
	ErrNoServerAddr       = errors.New("no server address")
	ErrNoNotificationAddr = errors.New("no notification address")
)

type Server struct {
	Addr       string
	NotifyAddr string
}

func (s *Server) load(envPrefix string) error {
	v := setupViper(envPrefix)

	s.Addr = v.GetString("address")
	if s.Addr == "" {
		return ErrNoServerAddr
	}

	s.NotifyAddr = v.GetString("notify.address")
	if s.NotifyAddr == "" {
		return ErrNoNotificationAddr
	}

	return nil
}
