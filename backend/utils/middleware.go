package utils

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Authorization header required",
		})
	}

	// Проверяем формат "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid authorization header format",
		})
	}

	token := parts[1]
	claims, err := ValidateToken(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Сохраняем данные пользователя в контексте
	c.Locals("userID", claims.UserID)
	c.Locals("userEmail", claims.Email)
	c.Locals("role", claims.Role)

	return c.Next()
}

// AdminMiddleware проверяет роль администратора
func AdminMiddleware(c *fiber.Ctx) error {
	roleVal := c.Locals("role")
	roleStr, ok := roleVal.(string)
	if !ok || roleStr != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Admin access required",
		})
	}
	return c.Next()
}
