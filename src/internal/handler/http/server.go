package http

import (
	"fmt"

	"pbmap_api/src/config"
	"pbmap_api/src/internal/repository"
	"pbmap_api/src/internal/usecase"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewServer(cfg *config.Config, handler *Handler) *fiber.App {
	app := fiber.New()

	api := app.Group("/api")
	api.Get("/health", handler.HealthCheck)

	// User routes
	users := api.Group("/users")
	users.Post("/", handler.UserHandler.Create)
	users.Get("/", handler.UserHandler.List)
	users.Get("/:id", handler.UserHandler.Get)
	users.Put("/:id", handler.UserHandler.Update)
	users.Delete("/:id", handler.UserHandler.Delete)

	return app
}

func Run(cfg *config.Config, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := NewUserHandler(userUsecase)

	handler := NewHandler(userHandler)
	app := NewServer(cfg, handler)

	addr := fmt.Sprintf(":%d", cfg.AppPort)
	if err := app.Listen(addr); err != nil {
		panic(err)
	}
}
