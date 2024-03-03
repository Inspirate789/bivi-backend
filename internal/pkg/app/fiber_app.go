package app

import (
	"context"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	_ "gitlab.teamdev.huds.su/bivi/backend/swagger"
	"log/slog"
)

type fiberApp struct {
	fiber  *fiber.App
	logger *slog.Logger
}

type ApiSettings struct {
	Port      string
	ApiPrefix string
}

func NewFiberApp(settings ApiSettings, log *slog.Logger) WebApp {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format: "${pid} ${locals:requestid} ${status} - ${latency} ${method} ${path}\n",
		Output: slog.NewLogLogger(log.Handler(), slog.LevelDebug).Writer(),
	}))

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:          fmt.Sprintf("http://localhost:%s/swagger/doc.json", settings.Port),
		DeepLinking:  false,
		DocExpansion: "none",
	}))

	// setup delivery

	return &fiberApp{
		fiber:  app,
		logger: log,
	}
}

func (f *fiberApp) Start(port string) error {
	return f.fiber.Listen(":" + port)
}

func (f *fiberApp) Stop(ctx context.Context) error {
	return f.fiber.ShutdownWithContext(ctx)
}
