package main

import (
	"app_burse_backend/configs"
	"app_burse_backend/internal/app/web"
)

// Запуск web-приложения.
func main() {
	config := configs.NewCondfig()
	conf := config.Load()

	webApp := web.NewWebContext(conf)
	webApp.InitDB()
	webApp.InitLogger()
	webApp.InitProducer()

	webApp.Run()
}
