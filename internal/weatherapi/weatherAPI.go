package weatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/enchik0reo/weatherTGBot/pkg/e"
)

const (
	requestPart1 = "https://api.openweathermap.org/data/2.5/weather?q="
	requestPart2 = "&units=metric&lang=ru&appid=79d1ca96933b0328e1c7e3e7a26cb347"
)

type Weather struct {
	Description string `json:"description"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
}

type Wind struct {
	Speed float64 `json:"speed"`
}

type WeatherForecast struct {
	Weather []Weather `json:"weather"`
	Main    Main      `json:"main"`
	Wind    Wind      `json:"wind"`
}

func GetWeatherForecast(city string) (*WeatherForecast, error) {
	req := fmt.Sprintf("%s%s%s", requestPart1, city, requestPart2)

	resp, err := http.Get(req)
	if err != nil {
		return nil, e.Wrap("can't get weather from api", err)
	}
	defer resp.Body.Close()

	wf := &WeatherForecast{}

	if err := json.NewDecoder(resp.Body).Decode(wf); err != nil {
		return nil, e.Wrap("can't decode response form api", err)
	}

	return wf, nil
}
