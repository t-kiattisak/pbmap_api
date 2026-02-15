package v1

import (
	"pbmap_api/src/internal/domain/entities"
	"pbmap_api/src/internal/dto"
	"pbmap_api/src/internal/usecase"
	"pbmap_api/src/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

// AlarmHandler handles alarm dispatch.
type AlarmHandler struct {
	alarmUsecase usecase.AlarmUsecase
	validator    *validator.Wrapper
}

// NewAlarmHandler creates the alarm HTTP handler.
func NewAlarmHandler(alarmUsecase usecase.AlarmUsecase, v *validator.Wrapper) *AlarmHandler {
	return &AlarmHandler{alarmUsecase: alarmUsecase, validator: v}
}

// Alarm dispatches an alarm (POST /api/v1/dispatch/alarm).
func (h *AlarmHandler) Alarm(c *fiber.Ctx) error {
	var req dto.AlarmDispatchRequest
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

	payload := &entities.AlarmDispatchRequest{
		AlarmID: req.AlarmID,
		Urgency: req.Urgency,
		Center:  entities.AlarmCenter{Lat: req.Center.Lat, Lng: req.Center.Lng, Radius: req.Center.Radius},
		Signal:  req.Signal,
		Content: req.Content,
	}

	if err := h.alarmUsecase.DispatchAlarm(c.Context(), payload); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Alarm dispatched successfully",
		Data:    map[string]string{"alarm_id": req.AlarmID},
	})
}
