package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/enchik0reo/weatherTelegramBot/internal/models"
	"github.com/enchik0reo/weatherTelegramBot/internal/weatherapi"
	"github.com/enchik0reo/weatherTelegramBot/pkg/e"
	"github.com/enchik0reo/weatherTelegramBot/pkg/mylogs"
)

const (
	RecentTime = 15 * time.Minute
)

type Storage interface {
	SaveWeatherHistory(f models.Forecast) error
	GetRecentForecasts() ([]models.Forecast, error)
	CloseConnect() error
}

type Cache interface {
	Save(forecast models.Forecast) error
	Show(city string) (models.Forecast, error)
	Exist(city string) bool
	Stop() error
}

type Repository struct {
	storage Storage
	cache   Cache
	log     *mylogs.Lgr
}

func New(s Storage, c Cache, log *mylogs.Lgr) (*Repository, error) {
	forecasts, err := s.GetRecentForecasts()
	if err != nil {
		return nil, e.Wrap("can't warmup cache", err)
	}

	for _, f := range forecasts {
		c.Save(f)
	}

	return &Repository{s, c, log}, nil
}

func (r *Repository) GetWeather(city string, userName string) (*models.Forecast, error) {
	f := models.Forecast{}

	if r.cache.Exist(city) {
		f, err := r.cache.Show(city)
		if err != nil {
			r.log.Debugf("can't get weather from cache %v", err)
		}
		return &f, nil
	}

	wf, err := weatherapi.GetWeatherForecast(city)
	if err != nil {
		if errors.Is(err, models.ErrCityNotFound) {
			return nil, err
		}
		return nil, e.Wrap("can't get weather", err)
	}

	f = models.Forecast{
		CityName:        city,
		UserName:        userName,
		ValidUntilUTC:   time.Now().UTC().Add(RecentTime),
		WeatherForecast: *wf,
	}

	r.cache.Save(f)

	err = r.storage.SaveWeatherHistory(f)
	if err != nil {
		err = e.Wrap("can't save to db", err)
	}

	return &f, err
}

func (r *Repository) CloseConnect() error {
	var msg string
	var err error

	err = r.cache.Stop()
	if err != nil {
		msg = fmt.Sprintf("%v", err)
	}

	err = r.storage.CloseConnect()
	if err != nil {
		if msg == "" {
			msg = fmt.Sprintf("%v", err)
		} else {
			msg = fmt.Sprintf("%s; %v", msg, err)
		}
	}

	if msg == "" {
		return nil
	} else {
		return e.Wrap("can't close repository connection", fmt.Errorf(msg))
	}
}
