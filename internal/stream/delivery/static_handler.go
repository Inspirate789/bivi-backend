package delivery

import (
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func StaticHandler(_ UseCase, logger *slog.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logger.Debug("file requested", slog.String("path", ctx.Path()))
		return ctx.Next()
	}
}
