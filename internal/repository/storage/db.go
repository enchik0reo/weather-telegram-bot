package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/enchik0reo/weatherTGBot/internal/models"
	"github.com/enchik0reo/weatherTGBot/pkg/e"

	_ "github.com/lib/pq"
)

type DBStorage struct {
	*sql.DB
}

func New(host, port, user, password, dbname, sslmode string) (*DBStorage, error) {
	connectStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}

	if err = db.Ping(); err != nil {
		return nil, e.Wrap("can't connect to database", err)
	}

	return &DBStorage{db}, nil
}

func (db *DBStorage) SaveWeatherHistory(f models.Forecast) error {
	q := `INSERT INTO weather VALUES ($1, $2, $3, $4)`

	jsonForecast, err := json.Marshal(f.WeatherForecast)
	if err != nil {
		return e.Wrap("can't save into db", err)
	}

	if _, err := db.Exec(q, f.CityName, f.UserName, f.ValidUntilUTC.Format("2006-01-02 15:04:05"), jsonForecast); err != nil {
		return e.Wrap("can't save into db", err)
	}

	return nil
}

func (db *DBStorage) GetRecentForecasts() ([]models.Forecast, error) {
	q := `SELECT city, user_name, valid_until, weather_forecast FROM weather WHERE valid_until > $1`

	qTime := time.Now().UTC().Format("2006-01-02 15:04:05")

	rows, err := db.Query(q, qTime)
	if err != nil {
		return nil, e.Wrap("can't load forecasts from db", err)
	}
	defer rows.Close()

	forecasts := []models.Forecast{}

	for rows.Next() {
		var forecast models.Forecast
		var jsonForecast []byte

		err = rows.Scan(&forecast.CityName, &forecast.UserName, &forecast.ValidUntilUTC, &jsonForecast)
		if err != nil {
			return nil, e.Wrap("can't scan forecast from db", err)
		}

		json.Unmarshal(jsonForecast, &forecast.WeatherForecast)

		forecasts = append(forecasts, forecast)
	}

	return forecasts, nil
}

func (db *DBStorage) CloseConnect() error {
	return db.Close()
}
