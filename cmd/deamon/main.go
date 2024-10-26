package main

import (
	"app_burse_backend/configs"
	"app_burse_backend/internal/app/deamon"
)

func main() {
	config := configs.NewCondfig()
	conf := config.Load()

	app := deamon.NewInstance(conf)
	app.Setup()
	app.Run()
}
