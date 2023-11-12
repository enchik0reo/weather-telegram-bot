package app

import (
	"github.com/enchik0reo/weatherTGBot/internal/clients/telegram"
	"github.com/enchik0reo/weatherTGBot/internal/config"
	"github.com/enchik0reo/weatherTGBot/internal/repository"
	"github.com/enchik0reo/weatherTGBot/internal/repository/cache"
	"github.com/enchik0reo/weatherTGBot/internal/repository/storage"
	"github.com/enchik0reo/weatherTGBot/pkg/mylogs"
)

const (
	tgBothost = "api.telegram.org"
)

type App struct {
	log      *mylogs.Lgr
	cfg      *config.Config
	tgClient *telegram.Client
}

func New() *App {
	var err error
	a := &App{}

	a.log = mylogs.New()

	a.cfg, err = config.Load()
	if err != nil {
		a.log.Fatalf("an error occurred while loading config: %v", err)
	}

	a.tgClient = telegram.New(tgBothost, a.cfg.Token)

	s, err := storage.New(a.cfg.DB.Host, a.cfg.DB.Port, a.cfg.DB.User, a.cfg.DB.User, a.cfg.DB.DBName, a.cfg.DB.SSLMode)
	if err != nil {
		a.log.Fatalf("an error occurred while creating new storage: %v", err)
	}

	c := cache.New()

	repo, err := repository.New(s, c)
	if err != nil {
		a.log.Fatalf("an error occurred while creating new repository: %v", err)
	}

	return a
}

func (a *App) Run() {

	// fetch := fetcher.New()

	// proces := processor.New()

	// consum := consumer.New(fetcher, processor)
	// consum.Start()
}
