package response

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
    Success bool        `json:"success"`
    Data    any         `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
    Meta    any         `json:"meta,omitempty"`
}

func JSON(c *fiber.Ctx, status int, data any) error {
    return c.Status(status).JSON(APIResponse{Success: true, Data: data})
}

func Error(c *fiber.Ctx, status int, err error) error {
    return c.Status(status).JSON(APIResponse{Success: false, Error: err.Error()})
}

func Message(c *fiber.Ctx, status int, message string) error {
    return c.Status(status).JSON(APIResponse{Success: true, Meta: map[string]string{"message": message}})
}
