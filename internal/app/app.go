package app

import (
	httpapp "effectivemobiletesttask/internal/app/http"
	"effectivemobiletesttask/internal/config"
	server "effectivemobiletesttask/internal/http-server/song"
	service "effectivemobiletesttask/internal/services/song"
	"effectivemobiletesttask/internal/storage/postgres"
	"effectivemobiletesttask/internal/utils/logger"
	"log/slog"
)

type App struct {
	HTTPserver *httpapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	storage, err := postgres.New(cfg.Storage)
	if err != nil {
		log.Error("failed connect to db err: %s", logger.Err(err))
	}

	service := service.New(log, storage, cfg.Client)
	server := server.New(log, cfg.PageSize, service)

	app := httpapp.New(log, &cfg.Server, server)

	return &App{
		HTTPserver: app,
	}
}
