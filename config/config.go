package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server Server
	Notify Notify
	DB     DB
}

func (c Config) Load() error {
	if err := c.Server.load("server"); err != nil {
		return fmt.Errorf("http server configuration: %w", err)
	}

	if err := c.Notify.load("notify"); err != nil {
		return fmt.Errorf("notified configuration: %w", err)
	}

	if err := c.DB.Load("storage"); err != nil {
		return fmt.Errorf("storage configuration: %w", err)
	}

	return nil
}

func setupViper(envPrefix string) *viper.Viper {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix(envPrefix)
	return v
}
