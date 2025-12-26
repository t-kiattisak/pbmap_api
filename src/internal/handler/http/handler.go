package http

import (
	"pbmap_api/src/domain"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	UserHandler *UserHandler
}

func NewHandler(userHandler *UserHandler) *Handler {
	return &Handler{
		UserHandler: userHandler,
	}
}

func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(domain.APIResponse{
		Status:  fiber.StatusOK,
		Message: "OK",
	})
}
