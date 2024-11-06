package main

import (
	"app_burse_backend/configs"
	"app_burse_backend/internal/app/web"
	"fmt"
)

// Запуск web-приложения.
func main() {
	config := configs.NewCondfig()
	conf := config.Load()

	fmt.Println("conf", conf)
	webApp := web.NewWebContext(conf)
	webApp.InitDB()
	webApp.InitLocales("")
	webApp.InitLogger()
	webApp.InitProducer()

	webApp.Run()
}
