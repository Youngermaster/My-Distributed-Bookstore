package response

import "github.com/gofiber/fiber/v2"

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// Success sends a success response
func Success(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(ErrorResponse{
		Success: false,
		Error:   message,
	})
}

// ErrorWithCode sends an error response with a code
func ErrorWithCode(c *fiber.Ctx, status int, message, code string) error {
	return c.Status(status).JSON(ErrorResponse{
		Success: false,
		Error:   message,
		Code:    code,
	})
}
