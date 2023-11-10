package app

import (
	"github.com/enchik0reo/weatherTGBot/internal/clients/telegram"
	"github.com/enchik0reo/weatherTGBot/pkg/mylogs"
)

type App struct {
	log      *mylogs.Lgr
	tgClient *telegram.Client
}

func New() *App {
	var err error
	a := &App{}

	a.log = mylogs.New()

	a.tgClient, err = telegram.New()
	if err != nil {
		a.log.Fatalf("an error occurred while creating the client: %v", err)
	}

	return a
}

func (a *App) Run() {

	// fetch := fetcher.New()

	// proces := processor.New()

	// consum := consumer.New(fetcher, processor)
	// consum.Start()
}
