package httpapp

import (
	"effectivemobiletesttask/internal/config"
	"effectivemobiletesttask/internal/http-server/song"
	"effectivemobiletesttask/internal/utils/logger"
	"fmt"
	"log/slog"
	"net/http"

	_ "effectivemobiletesttask/cmd/song-lib/docs"

	"github.com/rs/cors"
	swagger "github.com/swaggo/http-swagger"
)

type App struct {
	log        *slog.Logger
	cfg        *config.HTTPServer
	httpServer *http.Server
}

func New(log *slog.Logger, cfg *config.HTTPServer, server *song.Server) *App {
	mux := http.NewServeMux()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	})
	handler := corsHandler.Handler(mux)

	mux.HandleFunc("/swagger/", swagger.WrapHandler)
	server.RegisterRoutes(mux)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      handler,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &App{
		log:        log,
		cfg:        cfg,
		httpServer: httpServer,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "app.http.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("host", a.cfg.Host),
		slog.Int("port", a.cfg.Port),
	)

	log.Info("starting HTTP server")

	if err := a.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "app.http.Stop"

	a.log.With(slog.String("op", op)).Info("Stopping HTTP server", slog.Int("port", a.cfg.Port))

	if err := a.httpServer.Close(); err != nil {
		a.log.Error("error while closing HTTP server", logger.Err(err))
	}
}
