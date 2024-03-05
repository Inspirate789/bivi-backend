package app

import (
	"context"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2" // TODO: replace with "github.com/gofiber/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/pkg/errors"
	_ "gitlab.teamdev.huds.su/bivi/backend/swagger" // include generated swagger documentation
	"log/slog"
)

type FiberApp struct {
	fiber  *fiber.App
	logger *slog.Logger
}

type APISettings struct {
	Port      string
	APIPrefix string
}

func NewFiberApp(settings APISettings, log *slog.Logger) *FiberApp {
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

	return &FiberApp{
		fiber:  app,
		logger: log,
	}
}

func (f *FiberApp) Start(port string) error {
	return errors.Wrap(f.fiber.Listen(":"+port), "start application")
}

func (f *FiberApp) Shutdown(ctx context.Context) error {
	return errors.Wrap(f.fiber.ShutdownWithContext(ctx), "stop application")
}
