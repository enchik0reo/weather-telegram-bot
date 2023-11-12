package models

import (
	"errors"
	"time"
)

var ErrCityNotFound = errors.New("api didn't found the city")

type weather struct {
	Description string `json:"description"`
}

type main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
}

type wind struct {
	Speed float64 `json:"speed"`
}

type WeatherForecast struct {
	Weather []weather `json:"weather"`
	Main    main      `json:"main"`
	Wind    wind      `json:"wind"`
	Cod     string    `json:"cod"`
}

type Forecast struct {
	CityName        string
	UserName        string
	ValidUntilUTC   time.Time
	WeatherForecast WeatherForecast
}
