package delivery

import (
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"strings"
)

func StaticHandler(contentRoute string, _ Streamer, logger *slog.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		logger.Debug("file requested", slog.String("path", strings.TrimPrefix(ctx.Path(), contentRoute)))
		return ctx.Next()
	}
}
