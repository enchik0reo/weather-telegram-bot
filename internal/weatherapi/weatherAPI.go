package weatherapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/enchik0reo/weatherTGBot/internal/models"
	"github.com/enchik0reo/weatherTGBot/pkg/e"
)

const (
	requestPart1 = "https://api.openweathermap.org/data/2.5/weather?q="
	requestPart2 = "&units=metric&lang=ru&appid=79d1ca96933b0328e1c7e3e7a26cb347"
)

func GetWeatherForecast(city string) (*models.WeatherForecast, error) {
	req := fmt.Sprintf("%s%s%s", requestPart1, city, requestPart2)

	resp, err := http.Get(req)
	if err != nil {
		return nil, e.Wrap("can't get weather from api", err)
	}
	defer resp.Body.Close()

	wf := &models.WeatherForecast{}

	if err := json.NewDecoder(resp.Body).Decode(wf); err != nil {
		return nil, e.Wrap("can't decode response form api", err)
	}

	var icod int

	switch f := wf.Cod.(type) {
	case string:
		icod, err = strconv.Atoi(f)
		if err != nil {
			return nil, e.Wrap("can't convert cod from api", err)
		}
	case float64:
		icod = int(f)
	}

	if icod == http.StatusNotFound {
		return nil, models.ErrCityNotFound
	}

	return wf, nil
}
