package http

import (
	"fmt"

	"pbmap_api/src/config"
	"pbmap_api/src/internal/middleware"
	"pbmap_api/src/internal/repository"
	"pbmap_api/src/internal/usecase"
	"pbmap_api/src/pkg/auth"
	"pbmap_api/src/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewServer(cfg *config.Config, handler *Handler, jwtService *auth.JWTService) *fiber.App {
	app := fiber.New()

	api := app.Group("/api")
	api.Get("/health", handler.HealthCheck)

	// User routes
	users := api.Group("/users")
	users.Post("/", handler.UserHandler.Create)
	users.Get("/", handler.UserHandler.List)
	users.Get("/me", middleware.Protected(jwtService), handler.UserHandler.Me)
	users.Get("/:id", handler.UserHandler.Get)
	users.Put("/:id", handler.UserHandler.Update)
	users.Delete("/:id", handler.UserHandler.Delete)

	return app
}

func Run(cfg *config.Config, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	v := validator.New()
	jwtService := auth.NewJWTService(cfg.JWTSecret)
	userHandler := NewUserHandler(userUsecase, v, jwtService)

	handler := NewHandler(userHandler)
	app := NewServer(cfg, handler, jwtService)

	addr := fmt.Sprintf(":%d", cfg.AppPort)
	if err := app.Listen(addr); err != nil {
		panic(err)
	}
}
