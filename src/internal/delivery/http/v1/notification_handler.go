package v1

import (
	"pbmap_api/src/internal/domain/entities"
	"pbmap_api/src/internal/dto"
	"pbmap_api/src/internal/usecase"
	"pbmap_api/src/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// NotificationHandler handles notification endpoints.
type NotificationHandler struct {
	notificationUsecase usecase.NotificationUsecase
	validator           *validator.Wrapper
}

// NewNotificationHandler creates the notification HTTP handler.
func NewNotificationHandler(notificationUsecase usecase.NotificationUsecase, v *validator.Wrapper) *NotificationHandler {
	return &NotificationHandler{notificationUsecase: notificationUsecase, validator: v}
}

// Broadcast handles POST /api/notifications/broadcast.
func (h *NotificationHandler) Broadcast(c *fiber.Ctx) error {
	var req dto.BroadcastRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if errors := h.validator.Validate(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(entities.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Data:    errors,
		})
	}

	if err := h.notificationUsecase.Broadcast(c.Context(), req.Title, req.Body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Broadcast sent successfully",
	})
}

// Subscribe handles POST /api/notifications/subscribe.
func (h *NotificationHandler) Subscribe(c *fiber.Ctx) error {
	var req dto.SubscribeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if errors := h.validator.Validate(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(entities.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Data:    errors,
		})
	}

	result, err := h.notificationUsecase.SubscribeToTopic(c.Context(), req.Tokens, "all_devices")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Subscribed successfully",
		Data:    result,
	})
}

// Unsubscribe handles POST /api/notifications/unsubscribe.
func (h *NotificationHandler) Unsubscribe(c *fiber.Ctx) error {
	var req dto.SubscribeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if errors := h.validator.Validate(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(entities.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Data:    errors,
		})
	}

	result, err := h.notificationUsecase.UnsubscribeFromTopic(c.Context(), req.Tokens, "all_devices")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Unsubscribed successfully",
		Data:    result,
	})
}
