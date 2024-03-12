package app

import (
	"context"
	swagger "github.com/arsmn/fiber-swagger/v2" // replace with "github.com/gofiber/swagger" ?
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/pkg/errors"
	_ "gitlab.teamdev.huds.su/bivi/backend/swagger" // include generated swagger documentation
	"log/slog"
	"time"
)

type FiberApp struct {
	fiber  *fiber.App
	logger *slog.Logger
}

type APISettings struct {
	Port        string
	ContentPath string
}

func NewFiberApp(settings APISettings, log *slog.Logger) *FiberApp {
	app := fiber.New()
	staticApp := fiber.New()
	app.Mount("/", staticApp)

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New())

	app.Get("/metrics", monitor.New(monitor.Config{Title: "bivi metrics page"}))
	app.Use(pprof.New())

	app.Get("/swagger/*", swagger.New(swagger.Config{
		DeepLinking:  false,
		DocExpansion: "none",
	}))

	staticApp.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))
	staticApp.Static("/", settings.ContentPath, fiber.Static{
		Compress:      true,
		CacheDuration: -1 * time.Second,
	})

	_ = app.Group("/api/v1")
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
