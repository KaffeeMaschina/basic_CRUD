package apiserver

import (
	"github.com/KaffeeMaschina/basic_CRUD/internal/app/routes"
	"github.com/KaffeeMaschina/basic_CRUD/internal/app/storage"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"os"
)

type APIServer struct {
	Config   *Config
	Logger   *slog.Logger
	Database *storage.Database
}

// New creates new server instance with config, storage and logger
func New(config *Config, database *storage.Database) *APIServer {
	apiserver := &APIServer{
		Config:   config,
		Database: database,
	}
	apiserver.NewLogger(apiserver.Config.LogLevel)
	return apiserver
}

// Start creates new fiber.App, defines all available routes and start server
func (s *APIServer) Start() error {
	s.Logger.Debug("starting server")

	app := fiber.New()

	routes.PublicRoutes(app, s.Database, s.Logger)

	return app.Listen(s.Config.BindAddr)
}

// NewLogger creates new logger for APIServer instance with a logLevel from .env
func (s *APIServer) NewLogger(logLevel string) {
	var level slog.Level

	switch logLevel {
	case "debug", "":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	opts := &slog.HandlerOptions{Level: level}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	s.Logger = logger
}
