package middleware

import (
	"strings"

	"pbmap_api/src/domain"
	"pbmap_api/src/pkg/auth"

	"github.com/gofiber/fiber/v2"
)

func Protected(jwtService *auth.JWTService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.APIResponse{
				Status:  fiber.StatusUnauthorized,
				Message: "Missing authorization header",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.APIResponse{
				Status:  fiber.StatusUnauthorized,
				Message: "Invalid authorization format",
			})
		}

		tokenString := parts[1]
		tokenDetails, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(domain.APIResponse{
				Status:  fiber.StatusUnauthorized,
				Message: "Invalid or expired token",
			})
		}

		// Store user info in context locals
		c.Locals("user_id", tokenDetails.UserID)
		c.Locals("role", tokenDetails.Role)

		return c.Next()
	}
}
