package cache

import (
	"sync"
	"time"

	"github.com/enchik0reo/weatherTGBot/internal/repository"
)

type TimeCache struct {
	m map[string]repository.Forecast
	sync.RWMutex
}

func New() *TimeCache {
	return &TimeCache{
		m:       make(map[string]repository.Forecast),
		RWMutex: sync.RWMutex{},
	}
}

func (c *TimeCache) Save(forecast repository.Forecast) {
	c.Lock()
	c.m[forecast.CityName] = forecast
	c.Unlock()
}

func (c *TimeCache) Show(city string) repository.Forecast {
	c.RLock()
	defer c.RUnlock()
	return c.m[city]
}

func (c *TimeCache) IsExist(city string) bool {
	c.RLock()
	defer c.RUnlock()
	if f, ok := c.m[city]; ok {
		if f.ValidUntilUTC.Before(time.Now().UTC()) {
			return true
		}
		delete(c.m, city)
	}
	return false
}
