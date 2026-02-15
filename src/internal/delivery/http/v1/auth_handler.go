package v1

import (
	"pbmap_api/src/internal/domain/entities"
	"pbmap_api/src/internal/dto"
	"pbmap_api/src/internal/usecase"
	"pbmap_api/src/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthHandler handles auth endpoints.
type AuthHandler struct {
	authUsecase usecase.AuthService
	validator   *validator.Wrapper
}

// NewAuthHandler creates the auth HTTP handler.
func NewAuthHandler(authUsecase usecase.AuthService, v *validator.Wrapper) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase, validator: v}
}

// LoginWithSocial handles POST /api/auth/login.
func (h *AuthHandler) LoginWithSocial(c *fiber.Ctx) error {
	var req dto.SocialLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Cannot parse JSON",
		})
	}

	if errors := h.validator.Validate(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(entities.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Data:    errors,
		})
	}

	resp, err := h.authUsecase.LoginWithSocial(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(entities.APIResponse{
			Status:  fiber.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Login successfully",
		Data:    resp,
	})
}

// Logout handles POST /api/auth/logout.
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(entities.APIResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	if err := h.authUsecase.Logout(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Logout successfully",
	})
}

// RefreshToken handles POST /api/auth/refresh.
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(entities.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Cannot parse JSON",
		})
	}

	if errors := h.validator.Validate(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(entities.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Data:    errors,
		})
	}

	resp, err := h.authUsecase.RefreshToken(c.Context(), &req)
	if err != nil {
		status := fiber.StatusUnauthorized
		if err.Error() != "refresh token expired" && err.Error() != "invalid refresh token" {
			status = fiber.StatusInternalServerError
		}
		return c.Status(status).JSON(entities.APIResponse{
			Status:  status,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(entities.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Token refreshed successfully",
		Data:    resp,
	})
}
