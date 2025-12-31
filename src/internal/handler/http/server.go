package http

import (
	"fmt"

	"pbmap_api/src/config"
	"pbmap_api/src/internal/handler"
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

	// Notification routes
	notifications := api.Group("/notifications")
	notifications.Post("/broadcast", handler.NotificationHandler.Broadcast)
	notifications.Post("/subscribe", handler.NotificationHandler.Subscribe)
	notifications.Post("/unsubscribe", handler.NotificationHandler.Unsubscribe)

	return app
}

func Run(cfg *config.Config, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)

	fcmService, err := usecase.NewFCMService(cfg)
	if err != nil {
		fmt.Printf("Warning: Failed to initialize FCM Service: %v\n", err)
	}
	v := validator.New()
	notificationHandler := handler.NewNotificationHandler(fcmService, v)

	jwtService := auth.NewJWTService(cfg.JWTSecret)
	userHandler := handler.NewUserHandler(userUsecase, v, jwtService)

	handler := NewHandler(userHandler, notificationHandler)
	app := NewServer(cfg, handler, jwtService)

	addr := fmt.Sprintf(":%d", cfg.AppPort)
	if err := app.Listen(addr); err != nil {
		panic(err)
	}
}
