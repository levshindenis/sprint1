package main

import (
	"github.com/levshindenis/sprint1/internal/app/handlers"
	"github.com/levshindenis/sprint1/internal/app/routers"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var server handlers.HStorage
	server.Init()

	return http.ListenAndServe(server.GetStartSA(), routers.MyRouter(server))
}
