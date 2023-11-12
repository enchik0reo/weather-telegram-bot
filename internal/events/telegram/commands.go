package telegram

import (
	"errors"
	"fmt"
	"strings"

	"github.com/enchik0reo/weatherTGBot/internal/models"
	"github.com/enchik0reo/weatherTGBot/pkg/e"
)

const (
	StartCmd = "/start"
	HelpCmd  = "/help"
)

func (p *EventProcessor) doCmd(text string, chatID int, userName string) error {
	text = strings.TrimSpace(text)

	switch text {
	case StartCmd:
		return p.sendHello(chatID)
	case HelpCmd:
		return p.sendHelp(chatID)
	default:
		return p.sendWeather(chatID, text, userName)
	}
}

func (p *EventProcessor) sendWeather(chatID int, city, userName string) error {

	forecast, err := p.repository.GetWeather(city, userName)
	if err != nil {
		if errors.Is(err, models.ErrCityNotFound) {
			return p.client.SendMessage(chatID, msgUnknownCity)
		}
		return e.Wrap("can't do command \"get weather\"", err)
	}

	answer := makeAnswer(city, forecast)

	err = p.client.SendMessage(chatID, answer)
	if err != nil {
		return e.Wrap("can't do command \"get weather\"", err)
	}

	return nil
}

func (p *EventProcessor) sendHelp(chatID int) error {
	return p.client.SendMessage(chatID, msgHelp)
}

func (p *EventProcessor) sendHello(chatID int) error {
	return p.client.SendMessage(chatID, msgStart)
}

func makeAnswer(city string, f *models.Forecast) string {
	ansver := strings.Builder{}

	ansver.WriteString(fmt.Sprintf("Сейчас в городе %s %0.1f °C, ", city, f.WeatherForecast.Main.Temp))
	ansver.WriteString(fmt.Sprintf("%s.\n\n", makeDescription(f)))
	ansver.WriteString(fmt.Sprintf("Ощущается как %0.0f °C.\n\n", f.WeatherForecast.Main.FeelsLike))

	ws := f.WeatherForecast.Wind.Speed

	switch {
	case ws < 5:
		ansver.WriteString(slowWind)
	case ws < 15:
		ansver.WriteString(middleWind)
	case ws < 25:
		ansver.WriteString(strongWind)
	default:
		ansver.WriteString(storm)
	}

	return ansver.String()
}

func makeDescription(f *models.Forecast) string {
	desc := strings.Builder{}
	for i, w := range f.WeatherForecast.Weather {
		if i == len(f.WeatherForecast.Weather)-1 {
			desc.WriteString(w.Description)
		} else {
			desc.WriteString(fmt.Sprintf("%s; ", w.Description))
		}
	}

	return desc.String()
}
