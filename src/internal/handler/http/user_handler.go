package http

import (
	"pbmap_api/src/domain"
	"pbmap_api/src/internal/dto"
	"pbmap_api/src/internal/usecase"

	"pbmap_api/src/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	usecase   usecase.UserUsecase
	validator *validator.Wrapper
}

func NewUserHandler(usecase usecase.UserUsecase, v *validator.Wrapper) *UserHandler {
	return &UserHandler{
		usecase:   usecase,
		validator: v,
	}
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateUserRequest
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

	user := domain.User{
		Email:       &req.Email,
		DisplayName: req.DisplayName,
		Role:        req.Role,
	}

	if err := h.usecase.CreateUser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(domain.APIResponse{
		Status:  fiber.StatusCreated,
		Message: "User created successfully",
		Data:    user,
	})
}

func (h *UserHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid ID format",
		})
	}

	user, err := h.usecase.GetUser(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(domain.APIResponse{
			Status:  fiber.StatusNotFound,
			Message: "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.APIResponse{
		Status:  fiber.StatusOK,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid ID format",
		})
	}

	var user domain.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: err.Error(),
		})
	}

	user.ID = id
	if err := h.usecase.UpdateUser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.APIResponse{
		Status:  fiber.StatusOK,
		Message: "User updated successfully",
		Data:    user,
	})
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Invalid ID format",
		})
	}

	if err := h.usecase.DeleteUser(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.APIResponse{
		Status:  fiber.StatusOK,
		Message: "User deleted successfully",
	})
}

func (h *UserHandler) List(c *fiber.Ctx) error {
	users, err := h.usecase.ListUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}
