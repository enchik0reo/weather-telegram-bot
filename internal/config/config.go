package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/enchik0reo/weatherTelegramBot/pkg/e"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Token string `env:"TOKEN_TG_BOT"`
	DB    struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `env:"POSTGRES_PASSWORD"`
		DBName   string `yaml:"dbname"`
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"db"`
	Cache struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"cache"`
}

func Load() (*Config, error) {
	var c = Config{}
	var err error

	err = readYml(&c)
	if err != nil {
		return nil, err
	}

	err = loadEnv(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func readYml(c *Config) error {
	p := filepath.Join("config", "cfg.yml")
	if err := cleanenv.ReadConfig(p, c); err != nil {
		_, err = cleanenv.GetDescription(c, nil)
		return e.Wrap("an error occurred while getting config", err)
	}

	return nil
}

func loadEnv(c *Config) error {
	if err := godotenv.Load(); err != nil {
		return e.Wrap("an error occurred while loading environment variables", err)
	}

	c.Token = os.Getenv("TOKEN_TG_BOT")
	if c.Token == "" {
		return e.Wrap("token is not specified in environment variables", fmt.Errorf("just write it down"))
	}

	c.DB.Password = os.Getenv("POSTGRES_PASSWORD")
	if c.DB.Password == "" {
		return e.Wrap("db password is not specified in environment variables", fmt.Errorf("just write it down"))
	}

	return nil
}
