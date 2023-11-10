package config

import (
	"github.com/enchik0reo/weatherTGBot/pkg/e"
)

type Config struct {
	token string
}

func Load() (*Config, error) {
	c := &Config{}

	t, err := getToken()
	if err != nil {
		return nil, e.Wrap("an error occurred while getting token", err)
	}

	c.token = t

	return c, nil
}

func (c *Config) Token() string {
	return c.token
}
