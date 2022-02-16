package config

import "errors"

var ErrNoServerAddr = errors.New("no server address")

type Server struct {
	Addr string
}

func (s *Server) load(envPrefix string) error {
	v := setupViper(envPrefix)

	s.Addr = v.GetString("address")
	if s.Addr == "" {
		return ErrNoServerAddr
	}

	return nil
}
