package cache

import (
	"sync"
	"time"

	"github.com/enchik0reo/weatherTelegramBot/internal/models"
)

type TimeCache struct {
	m map[string]models.Forecast
	sync.RWMutex
}

func NewMapCache() *TimeCache {
	return &TimeCache{
		m:       make(map[string]models.Forecast),
		RWMutex: sync.RWMutex{},
	}
}

func (c *TimeCache) Save(forecast models.Forecast) {
	c.Lock()
	c.m[forecast.CityName] = forecast
	c.Unlock()
}

func (c *TimeCache) Show(city string) models.Forecast {
	c.RLock()
	defer c.RUnlock()
	return c.m[city]
}

func (c *TimeCache) IsExist(city string) bool {
	c.RLock()
	defer c.RUnlock()
	if f, ok := c.m[city]; ok {
		if time.Now().UTC().Before(f.ValidUntilUTC) {
			return true
		}
		delete(c.m, city)
	}
	return false
}
