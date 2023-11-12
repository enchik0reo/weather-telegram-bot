package telegram

import (
	"github.com/enchik0reo/weatherTGBot/internal/clients/telegram"
	"github.com/enchik0reo/weatherTGBot/internal/repository"
)

type Repository interface {
	GetWeather(city string, userName string) (*repository.Forecast, error)
}

type EventProcessor struct {
	tg     *telegram.Client
	offset int
	repo   Repository
}

func New(client *telegram.Client, repository Repository) *EventProcessor {

}
