package middleware

import (
	"strings"

	"github.com/Youngermaster/My-Distributed-Bookstore/admin-service/pkg/jwt"
	"github.com/Youngermaster/My-Distributed-Bookstore/admin-service/pkg/response"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Error(c, fiber.StatusUnauthorized, "Missing authorization header")
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid authorization header format")
		}

		tokenString := parts[1]

		// Validate token
		claims, err := jwt.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "Invalid or expired token")
		}

		// Store claims in context for use in handlers
		c.Locals("user_id", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("roles", claims.Roles)
		c.Locals("claims", claims)

		return c.Next()
	}
}

// RequireRole middleware ensures user has specific role
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*jwt.Claims)
		if !ok {
			return response.Error(c, fiber.StatusUnauthorized, "Unauthorized")
		}

		// Check if user has any of the required roles
		if !claims.HasAnyRole(roles...) {
			return response.Error(c, fiber.StatusForbidden, "Insufficient permissions")
		}

		return c.Next()
	}
}

// RequireAdmin middleware ensures user is admin
func RequireAdmin() fiber.Handler {
	return RequireRole("admin")
}
