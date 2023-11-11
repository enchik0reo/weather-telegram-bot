package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/enchik0reo/weatherTGBot/internal/repository"
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

func (db *DBStorage) SaveWeatherHistory(f repository.Forecast) error {
	q := `INSERT INTO weather VALUES ($1, $2, $3, $4)`

	data, err := json.MarshalIndent(f.WeatherForecast, "", " ")
	if err != nil {
		return e.Wrap("can't save into db", err)
	}

	textData := fmt.Sprintf("%x", data)

	if _, err := db.Exec(q, f.CityName, f.UserName, f.ValidUntilUTC, textData); err != nil {
		return e.Wrap("can't save into db", err)
	}

	return nil
}

func (db *DBStorage) GetRecentForecasts() ([]repository.Forecast, error) {
	q := `SELECT forecast FROM weather WHERE valid_until_utc > $1` // How to make this work?

	rows, err := db.Query(q, time.Now().UTC()) // Huh?
	if err != nil {
		return nil, e.Wrap("can't load forecasts from db", err)
	}
	defer rows.Close()

	forecasts := []repository.Forecast{}

	for rows.Next() {
		var forecast repository.Forecast
		var textData string

		err = rows.Scan(&forecast.CityName, &forecast.UserName, &forecast.ValidUntilUTC, &textData)
		if err != nil {
			return nil, e.Wrap("can't scan forecast from db", err)
		}

		json.Unmarshal([]byte(textData), &forecast.WeatherForecast)

		forecasts = append(forecasts, forecast)
	}

	return forecasts, nil
}

func (db *DBStorage) CloseConnect() error {
	return db.Close()
}
