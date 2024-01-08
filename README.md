# # Weather Telegram Bot

It's a Go application called Weather Telegram Bot which allows people to get weather forecast for any city by using telegram.

## Features

- Application uses weather API for forecast
- User can write name of any city to bot
- Bot takes the name of the city and returns the weather in it
- Application uses standart net/http pacage
- JSON communication
- Data persistence using PostgreSQL

## Development

Software requirements:

- Go
- Docker

To start the application:

```sh
$ git clone https://github.com/enchik0reo/weatherTelegramBot
$ cd weatherTelegramBot

# Create your telegram bot and get token

# Create an .env file with your environment variables: TOKEN_TG_BOT and POSTGRES_PASSWORD

# Run app
$ docker-compose up --build go
```

To terminate the application use `SIGTERM` signal (use Ctrl+C).