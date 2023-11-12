package events

import "github.com/enchik0reo/weatherTGBot/internal/models"

type Fetcher interface {
	Fetch(limit int) ([]models.Event, error)
}

type Processor interface {
	Process(e models.Event) error
}
