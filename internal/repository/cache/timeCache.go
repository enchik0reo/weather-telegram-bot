package cache

import (
	"time"

	"github.com/enchik0reo/weatherTGBot/internal/repository"
)

type TimeCache struct {
	m map[string]repository.Forecast
}

func New() *TimeCache {
	return &TimeCache{make(map[string]repository.Forecast)}
}

func (c *TimeCache) Save(forecast repository.Forecast) {
	c.m[forecast.CityName] = forecast
}

func (c *TimeCache) Show(city string) repository.Forecast {
	return c.m[city]
}

func (c *TimeCache) IsExist(city string) bool {
	if f, ok := c.m[city]; ok {
		if f.CreatedAt.Before(time.Now().Add(repository.RecentTime)) {
			return true
		}
		delete(c.m, city)
	}
	return false
}
