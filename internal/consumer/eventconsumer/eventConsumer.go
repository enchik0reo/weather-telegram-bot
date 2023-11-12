package eventconsumer

import (
	"sync"
	"time"

	"github.com/enchik0reo/weatherTGBot/internal/models"
	"github.com/enchik0reo/weatherTGBot/pkg/mylogs"
)

type Fetcher interface {
	Fetch(limit int) ([]models.Event, error)
}

type Processor interface {
	Process(e models.Event) error
}

type Consumer struct {
	fetcher   Fetcher
	processor Processor
	log       *mylogs.Lgr
	batchSize int
}

func New(f Fetcher, p Processor, l *mylogs.Lgr, bs int) *Consumer {
	return &Consumer{f, p, l, bs}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			c.log.Errorf("consumer: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			c.log.Error(err)
			continue
		}
	}
}

func (c *Consumer) handleEvents(events []models.Event) error {
	wg := sync.WaitGroup{}

	for _, event := range events {
		wg.Add(1)

		go func(event models.Event, wg *sync.WaitGroup) {
			defer wg.Done()

			c.log.Infof("got new event: %s", event.Text)

			if err := c.processor.Process(event); err != nil {
				c.log.Errorf("can't handle event: %s", err.Error())
			}
		}(event, &wg)
	}
	wg.Wait()

	return nil
}
