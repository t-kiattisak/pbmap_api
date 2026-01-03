package handler

import (
	"pbmap_api/src/domain"
	"pbmap_api/src/internal/dto"
	"pbmap_api/src/internal/usecase"
	"pbmap_api/src/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService usecase.AuthService
	validator   *validator.Wrapper
}

func NewAuthHandler(authService usecase.AuthService, v *validator.Wrapper) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   v,
	}
}

func (h *AuthHandler) LoginWithSocial(c *fiber.Ctx) error {
	var req dto.SocialLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Cannot parse JSON",
		})
	}

	if errors := h.validator.Validate(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Data:    errors,
		})
	}

	resp, err := h.authService.LoginWithSocial(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.APIResponse{
			Status:  fiber.StatusUnauthorized,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Login successfully",
		Data:    resp,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(domain.APIResponse{
			Status:  fiber.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	if err := h.authService.Logout(c.Context(), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(domain.APIResponse{
			Status:  fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Logout successfully",
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Cannot parse JSON",
		})
	}

	if errors := h.validator.Validate(req); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(domain.APIResponse{
			Status:  fiber.StatusBadRequest,
			Message: "Validation failed",
			Data:    errors,
		})
	}

	resp, err := h.authService.RefreshToken(c.Context(), &req)
	if err != nil {
		status := fiber.StatusUnauthorized
		if err.Error() == "refresh token expired" || err.Error() == "invalid refresh token" {
			status = fiber.StatusUnauthorized
		} else {
			status = fiber.StatusInternalServerError
		}

		return c.Status(status).JSON(domain.APIResponse{
			Status:  status,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(domain.APIResponse{
		Status:  fiber.StatusOK,
		Message: "Token refreshed successfully",
		Data:    resp,
	})
}
