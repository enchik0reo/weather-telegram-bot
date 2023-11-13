package main

import "github.com/enchik0reo/weatherTGBot/internal/app"

// ev.Text и meta.UserName почему то одиноковые и несут в себе city

func main() {
	app.New().Run()
}
