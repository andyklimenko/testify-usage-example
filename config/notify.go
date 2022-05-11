package config

import "errors"

var ErrNoNotificationAddr = errors.New("no notification address")

type Notify struct {
	Addr string
}

func (n *Notify) load(envPrefix string) error {
	v := setupViper(envPrefix)

	n.Addr = v.GetString("address")
	if n.Addr == "" {
		return ErrNoNotificationAddr
	}

	return nil
}
