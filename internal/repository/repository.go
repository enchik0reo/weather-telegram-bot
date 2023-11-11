package repository

import (
	"time"

	"github.com/enchik0reo/weatherTGBot/internal/weatherapi"
	"github.com/enchik0reo/weatherTGBot/pkg/e"
)

const (
	RecentTime = 30 * time.Minute
)

type Storage interface {
	SaveWeatherHistory(forecast Forecast) error
	GetRecentForecasts(time.Duration) ([]Forecast, error)
}

type Cache interface {
	Save(forecast Forecast)
	Show(city string) Forecast
	IsExist(city string) bool
}

type Repository struct {
	storage Storage
	cache   Cache
}

func New(s Storage, c Cache) (*Repository, error) {
	forecasts, err := s.GetRecentForecasts(RecentTime)
	if err != nil {
		return nil, e.Wrap("can't warmup cache", err)
	}

	for _, f := range forecasts {
		c.Save(f)
	}

	return &Repository{s, c}, nil
}

func (r Repository) SaveWeather(forecast Forecast) error {
	r.cache.Save(forecast)

	err := r.storage.SaveWeatherHistory(forecast)
	if err != nil {
		err = e.Wrap("can't save to db", err)
	}

	return err
}

func (r Repository) GetWeather(city string) Forecast {
	return r.cache.Show(city)
}

func (r Repository) IsExist(city string) bool {
	return r.cache.IsExist(city)
}

type Forecast struct {
	CityName        string
	UserName        string
	CreatedAt       time.Time
	WeatherForecast weatherapi.WeatherForecast
}
