package main

import (
	"effectivemobiletesttask/internal/app"
	"effectivemobiletesttask/internal/config"
	"effectivemobiletesttask/internal/migrator"
	"effectivemobiletesttask/internal/utils/logger"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

// @title Swagger Song Lib API
// @version 1.0
// @termsOfService http://swagger.io/terms/

// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8000
// @BasePath /
func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	log.Info("start 'song library' application")

	migrator.RunMigrations(&cfg.Storage, &cfg.Migrations)

	application := app.New(log, cfg)

	application.HTTPserver.MustRun()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping application", slog.String("signal:", sign.String()))

	application.HTTPserver.Stop()

	log.Info("server is dead")
}
