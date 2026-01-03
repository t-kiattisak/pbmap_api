package http

import (
	"fmt"

	"pbmap_api/src/config"
	"pbmap_api/src/internal/handler"
	"pbmap_api/src/internal/middleware"
	"pbmap_api/src/internal/repository"
	"pbmap_api/src/internal/usecase"
	"pbmap_api/src/pkg/auth"
	"pbmap_api/src/pkg/redis"
	"pbmap_api/src/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewServer(cfg *config.Config, handler *Handler, jwtService *auth.JWTService, tokenRepo repository.TokenRepository) *fiber.App {
	app := fiber.New()

	api := app.Group("/api")
	api.Get("/health", handler.HealthCheck)

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", handler.AuthHandler.LoginWithSocial)
	auth.Post("/logout", middleware.Protected(jwtService, tokenRepo), handler.AuthHandler.Logout)
	auth.Post("/refresh", handler.AuthHandler.RefreshToken)

	// User routes
	users := api.Group("/users")
	users.Post("/", handler.UserHandler.Create)
	users.Get("/", handler.UserHandler.List)
	users.Get("/me", middleware.Protected(jwtService, tokenRepo), handler.UserHandler.Me)
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
	deviceRepo := repository.NewDeviceRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, deviceRepo)

	fcmService, err := usecase.NewFCMService(cfg)
	if err != nil {
		fmt.Printf("Warning: Failed to initialize FCM Service: %v\n", err)
	}

	redisClient, err := redis.NewRedisClient(cfg)
	if err != nil {
		fmt.Printf("Warning: Failed to connect to Redis: %v\n", err)
	} else {
		fmt.Println("Successfully connected to Redis")
	}

	v := validator.New()
	notificationHandler := handler.NewNotificationHandler(fcmService, v)

	tokenRepo := repository.NewTokenRepository(redisClient)

	jwtService := auth.NewJWTService(cfg.JWTSecret)

	sessionRepo := repository.NewSessionRepository(db)
	tm := repository.NewTransactionManager(db)

	authService := usecase.NewAuthService(userUsecase, tokenRepo, sessionRepo, tm, jwtService, cfg)
	authHandler := handler.NewAuthHandler(authService, v)

	userHandler := handler.NewUserHandler(userUsecase, v, jwtService)

	handler := NewHandler(userHandler, notificationHandler, authHandler)
	app := NewServer(cfg, handler, jwtService, tokenRepo)

	addr := fmt.Sprintf(":%d", cfg.AppPort)
	if err := app.Listen(addr); err != nil {
		panic(err)
	}
}
