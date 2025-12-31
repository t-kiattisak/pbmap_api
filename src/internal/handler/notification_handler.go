package handler

import (
	"pbmap_api/src/domain"
	"pbmap_api/src/internal/dto"
	"pbmap_api/src/internal/usecase"
	"pbmap_api/src/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type NotificationHandler struct {
	fcmService usecase.FCMService
	validator  *validator.Wrapper
}

func NewNotificationHandler(fcmService usecase.FCMService, v *validator.Wrapper) *NotificationHandler {
	return &NotificationHandler{
		fcmService: fcmService,
		validator:  v,
	}
}

func (h *NotificationHandler) Broadcast(c *fiber.Ctx) error {
	var req dto.BroadcastRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if errors := h.validator.Validate(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Data:    errors,
		})
	}

	if err := h.fcmService.BroadcastNotification(c.Context(), req.Title, req.Body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Broadcast sent successfully",
	})
}

func (h *NotificationHandler) Subscribe(c *fiber.Ctx) error {
	var req dto.SubscribeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if errors := h.validator.Validate(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Data:    errors,
		})
	}

	result, err := h.fcmService.SubscribeToTopic(c.Context(), req.Tokens, "all_devices")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Subscribed successfully",
		Data:    result,
	})
}

func (h *NotificationHandler) Unsubscribe(c *fiber.Ctx) error {
	var req dto.SubscribeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if errors := h.validator.Validate(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Data:    errors,
		})
	}

	result, err := h.fcmService.UnsubscribeFromTopic(c.Context(), req.Tokens, "all_devices")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Unsubscribed successfully",
		Data:    result,
	})
}
