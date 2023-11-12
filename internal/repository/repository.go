package repository

import (
	"time"

	"github.com/enchik0reo/weatherTGBot/internal/weatherapi"
	"github.com/enchik0reo/weatherTGBot/pkg/e"
)

const (
	RecentTime = 15 * time.Minute
)

type Storage interface {
	SaveWeatherHistory(f Forecast) error
	GetRecentForecasts() ([]Forecast, error)
	CloseConnect() error
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

type Forecast struct {
	CityName        string
	UserName        string
	ValidUntilUTC   time.Time
	WeatherForecast weatherapi.WeatherForecast
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

func (r *Repository) GetWeather(city string, userName string) (*Forecast, error) {
	if r.cache.IsExist(city) {
		f := r.cache.Show(city)
		return &f, nil
	}

	wf, err := weatherapi.GetWeatherForecast(city)
	if err != nil {
		return nil, e.Wrap("can't get weather", err)
	}

	f := Forecast{
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
