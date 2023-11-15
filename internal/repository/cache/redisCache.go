package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/enchik0reo/weatherTelegramBot/internal/models"
	"github.com/enchik0reo/weatherTelegramBot/pkg/e"

	"github.com/gomodule/redigo/redis"
)

type RedisCache struct {
	c redis.Conn
}

func NewRedis(host, port string) (*RedisCache, error) {
	r := RedisCache{}
	conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return nil, e.Wrap("can't connect to redis", err)
	}
	r.c = conn

	return &r, nil
}

func (r *RedisCache) Save(forecast models.Forecast) error {
	if !r.Exist(forecast.CityName) {
		jsn, err := json.Marshal(forecast)
		if err != nil {
			return e.Wrap("can't marshal forecast", err)
		}

		_, err = r.c.Do("APPEND", forecast.CityName, jsn)
		if err != nil {
			return e.Wrap("can't append forecast to cache", err)
		}
	}

	return nil
}

func (r *RedisCache) Show(city string) (models.Forecast, error) {
	f := models.Forecast{}

	res, err := redis.String(r.c.Do("GET", city))
	if err != nil {
		return f, e.Wrap("can't get weather form cache", err)
	}

	err = json.Unmarshal([]byte(res), &f)
	if err != nil {
		return f, e.Wrap("can't unmarshal forecast", err)
	}

	return f, nil
}

func (r *RedisCache) Exist(city string) bool {
	f := models.Forecast{}

	res, err := redis.String(r.c.Do("GET", city))
	if err != nil {
		return false
	}

	err = json.Unmarshal([]byte(res), &f)
	if err != nil {
		return false
	}

	if time.Now().UTC().Before(f.ValidUntilUTC) {
		return true
	}

	_, err = r.c.Do("del", city)
	if err != nil {
		return false
	}

	return false
}

func (r *RedisCache) Stop() error {
	return r.c.Close()
}
