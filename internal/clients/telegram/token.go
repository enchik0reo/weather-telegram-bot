package telegram

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func mustToken() (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", fmt.Errorf("an error occurred while loading environment variables: %v", err)
	}

	token := os.Getenv("TOKEN_TG_BOT")

	if token == "" {
		return "", fmt.Errorf("the token is not specified in environment variables")
	}

	return token, nil
}
