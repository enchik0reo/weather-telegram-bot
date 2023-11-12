package repository

import (
	"errors"
	"time"

	"github.com/enchik0reo/weatherTGBot/internal/models"
	"github.com/enchik0reo/weatherTGBot/internal/weatherapi"
	"github.com/enchik0reo/weatherTGBot/pkg/e"
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
	Save(forecast models.Forecast)
	Show(city string) models.Forecast
	IsExist(city string) bool
}

type Repository struct {
	storage Storage
	cache   Cache
}

func New(s Storage, c Cache) (*Repository, error) {
	forecasts, err := s.GetRecentForecasts()
	if err != nil {
		return nil, e.Wrap("can't warmup cache", err)
	}

	for _, f := range forecasts {
		c.Save(f)
	}

	return &Repository{s, c}, nil
}

func (r *Repository) GetWeather(city string, userName string) (*models.Forecast, error) {
	if r.cache.IsExist(city) {
		f := r.cache.Show(city)
		return &f, nil
	}

	wf, err := weatherapi.GetWeatherForecast(city)
	if err != nil {
		if errors.Is(err, models.ErrCityNotFound) {
			return nil, err
		}
		return nil, e.Wrap("can't get weather", err)
	}

	f := models.Forecast{
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
	return r.storage.CloseConnect()
}
