package main

import (
	"log/slog"
	"os"

	"github.com/Xapsiel/EffectiveMobile/internal/config"
	"github.com/Xapsiel/EffectiveMobile/internal/handler"
	"github.com/Xapsiel/EffectiveMobile/internal/model"
	"github.com/Xapsiel/EffectiveMobile/internal/repository"
	"github.com/Xapsiel/EffectiveMobile/internal/service"
)

//	@title			Songs API
// 	@version 1.0
//	@description	This is an API for managing songs.

// @host			localhost:8080
// @BasePath		/
func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	db, err := repository.NewPostgresDB(cfg.DatabaseConfig)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(model.Server)
	if err := srv.Run(cfg.HostConfig.Port, handlers.InitRoutes()); err != nil {
		slog.Error(err.Error())
		os.Exit(1)

	}
}
