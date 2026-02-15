package v1

import (
	"pbmap_api/src/internal/domain/entities"
	"pbmap_api/src/internal/dto"
	"pbmap_api/src/internal/usecase"
	"pbmap_api/src/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PotentialPointHandler struct {
	usecase   usecase.PotentialPointUsecase
	validator *validator.Wrapper
}

// NewPotentialPointHandler creates the handler.
func NewPotentialPointHandler(usecase usecase.PotentialPointUsecase, v *validator.Wrapper) *PotentialPointHandler {
	return &PotentialPointHandler{
		usecase:   usecase,
		validator: v,
	}
}

// Create handles POST /api/v1/potential-points
func (h *PotentialPointHandler) Create(c *fiber.Ctx) error {
	var req dto.CreatePotentialPointInput
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

	var creatorID *uuid.UUID
	if userID, ok := c.Locals("user_id").(uuid.UUID); ok {
		creatorID = &userID
	}

	pp, err := h.usecase.Create(c.Context(), req, creatorID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(entities.APIResponse{
		Status:  fiber.StatusCreated,
		Message: "Potential point created successfully",
		Data:    pp,
	})
}

// Get handles GET /api/v1/potential-points/:id
func (h *PotentialPointHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid ID format",
		})
	}

	pp, err := h.usecase.FindByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(entities.APIResponse{
			Status:  fiber.StatusNotFound,
			Message: "Potential point not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Potential point retrieved successfully",
		Data:    pp,
	})
}

// Update handles PUT /api/v1/potential-points/:id
func (h *PotentialPointHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid ID format",
		})
	}

	var req dto.UpdatePotentialPointInput
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	pp, err := h.usecase.Update(c.Context(), id, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Potential point updated successfully",
		Data:    pp,
	})
}

// Delete handles DELETE /api/v1/potential-points/:id
func (h *PotentialPointHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid ID format",
		})
	}

	if err := h.usecase.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Potential point deleted successfully",
	})
}

// List handles GET /api/v1/potential-points
func (h *PotentialPointHandler) List(c *fiber.Ctx) error {
	pps, err := h.usecase.FindAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Potential points retrieved successfully",
		Data:    pps,
	})
}
