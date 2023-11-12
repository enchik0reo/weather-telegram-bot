package app

import (
	tgClient "github.com/enchik0reo/weatherTGBot/internal/clients/telegram"
	"github.com/enchik0reo/weatherTGBot/internal/config"
	"github.com/enchik0reo/weatherTGBot/internal/consumer/eventconsumer"
	"github.com/enchik0reo/weatherTGBot/internal/events/telegram"
	"github.com/enchik0reo/weatherTGBot/internal/repository"
	"github.com/enchik0reo/weatherTGBot/internal/repository/cache"
	"github.com/enchik0reo/weatherTGBot/internal/repository/storage"
	"github.com/enchik0reo/weatherTGBot/pkg/mylogs"
)

const (
	tgBothost = "api.telegram.org"
	batchSize = 100
)

type App struct {
	log           *mylogs.Lgr
	cfg           *config.Config
	tgClient      *tgClient.Client
	tgEventProc   *telegram.EventProcessor
	eventConsumer *eventconsumer.Consumer
}

func New() *App {
	var err error
	a := &App{}

	a.log = mylogs.New()

	a.cfg, err = config.Load()
	if err != nil {
		a.log.Fatalf("an error occurred while loading config: %v", err)
	}

	a.tgClient = tgClient.New(tgBothost, a.cfg.Token)

	s, err := storage.New(a.cfg.DB.Host, a.cfg.DB.Port, a.cfg.DB.User, a.cfg.DB.Password, a.cfg.DB.DBName, a.cfg.DB.SSLMode)
	if err != nil {
		a.log.Fatalf("an error occurred while creating new storage: %v", err)
	}

	c := cache.New()

	repo, err := repository.New(s, c)
	if err != nil {
		a.log.Fatalf("an error occurred while creating new repository: %v", err)
	}

	a.tgEventProc = telegram.New(a.tgClient, repo)

	a.eventConsumer = eventconsumer.New(a.tgEventProc, a.tgEventProc, a.log, batchSize)

	return a
}

func (a *App) Run() {
	if err := a.eventConsumer.Start(); err != nil {
		a.log.Fatalf("the consumer stoped %v", err)
	}
}
