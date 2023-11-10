package config

import (
	"fmt"
	"os"

	"github.com/enchik0reo/weatherTGBot/pkg/e"

	"github.com/joho/godotenv"
)

func getToken() (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", e.Wrap("an error occurred while loading environment variables", err)
	}

	token := os.Getenv("TOKEN_TG_BOT")

	if token == "" {
		return "", e.Wrap("token is not specified in environment variables", fmt.Errorf("just write it down"))
	}

	return token, nil
}
