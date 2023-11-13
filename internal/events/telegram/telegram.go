package telegram

import (
	"errors"

	"github.com/enchik0reo/weatherTGBot/internal/models"
	"github.com/enchik0reo/weatherTGBot/pkg/e"
)

type Repository interface {
	GetWeather(city string, userName string) (*models.Forecast, error)
	CloseConnect() error
}

type Client interface {
	GetUpdates(offset int, limit int) ([]models.Update, error)
	SendMessage(chatId int, text string) error
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

type Meta struct {
	ChatID   int
	UserName string
}

type EventProcessor struct {
	client     Client
	offset     int
	repository Repository
}

func New(c Client, r Repository) *EventProcessor {
	return &EventProcessor{client: c, repository: r}
}

func (p *EventProcessor) Fetch(limit int) ([]models.Event, error) {
	updates, err := p.client.GetUpdates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	events := make([]models.Event, 0, len(updates))

	for _, u := range updates {
		events = append(events, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return events, nil
}

func (p *EventProcessor) Process(ev models.Event) error {
	switch ev.Type {
	case models.Message:
		return p.processMessage(ev)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *EventProcessor) Stop() error {
	return p.repository.CloseConnect()
}

func (p *EventProcessor) processMessage(ev models.Event) error {
	meta, err := meta(ev)
	if err != nil {
		return e.Wrap("can't get process message", err)
	}

	if err = p.doCmd(ev.Text, meta.ChatID, meta.UserName); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(ev models.Event) (Meta, error) {
	met, ok := ev.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return met, nil
}

func event(upd models.Update) models.Event {
	updType := fetchType(upd)
	ev := models.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	if updType == models.Message {
		ev.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserName: upd.Message.From.UserName,
		}
	}

	return ev
}

func fetchType(upd models.Update) models.Type {
	if upd.Message == nil {
		return models.Unknown
	}
	return models.Message
}

func fetchText(upd models.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}
