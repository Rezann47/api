package app

import (
	"log"

	"go-gin-crud/config"
	"go-gin-crud/internal/server"
)

type App struct {
	cfg *config.Config
}

func New() *App {
	cfg := config.Load()
	return &App{cfg: cfg}
}

func (a *App) Run() {
	srv := server.NewHTTPServer(a.cfg)

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
