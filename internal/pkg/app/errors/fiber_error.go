package errors

import "github.com/gofiber/fiber/v2"

type FiberError struct {
	Error string `json:"error"`
}

func NewFiberError(msg string) fiber.Map {
	return fiber.Map{"error": msg}
}
