package telegram

import (
	"time"

	"github.com/enchik0reo/weatherTGBot/internal/clients/telegram"
)

const (
	RecentTime = 15 * time.Minute
)

type EventProcessor struct {
	tg     *telegram.Client
	offset int
	// weatherAPI
	// repository
}

func New(client *telegram.Client) {

}
