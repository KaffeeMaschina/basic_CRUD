package routes

import (
	"github.com/KaffeeMaschina/basic_CRUD/internal/app/storage"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

// PublicRoutes defines all available routes
func PublicRoutes(a *fiber.App, database *storage.Database, logger *slog.Logger) {
	route := a.Group("/api/v1")
	store := storage.Storage{DB: database, Logger: logger}
	route.Post("/tasks", store.CreateTask)
	route.Get("/tasks", store.GetTasks)
	route.Put("/tasks/:id", store.UpdateTask)
	route.Delete("/tasks/:id", store.DeleteTask)
}
