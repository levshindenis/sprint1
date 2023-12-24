// пакеты исполняемых приложений должны называться main
package main

import (
	"github.com/levshindenis/sprint1/internal/app/routers"
	"net/http"

	"github.com/levshindenis/sprint1/internal/app/handlers"
)

// функция main вызывается автоматически при запуске приложения
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

// функция run будет полезна при инициализации зависимостей сервера перед запуском
func run() error {
	var server handlers.HStorage
	server.Init()

	return http.ListenAndServe(server.GetStartSA(), routers.MyRouter(server))
}
