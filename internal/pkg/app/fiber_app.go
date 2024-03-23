package app

import (
	"context"
	swagger "github.com/arsmn/fiber-swagger/v2" // replace with "github.com/gofiber/swagger" ?
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/pkg/errors"
	clientDelivery "gitlab.teamdev.huds.su/bivi/backend/internal/client/delivery"
	streamDelivery "gitlab.teamdev.huds.su/bivi/backend/internal/stream/delivery"
	_ "gitlab.teamdev.huds.su/bivi/backend/swagger" // include generated swagger documentation
	"log/slog"
	"time"
)

type FiberApp struct {
	fiber  *fiber.App
	logger *slog.Logger
}

type APISettings struct {
	Port                string
	ContentPath         string
	ClientLogPath       string
	UploadFilesizeLimit int64
}

func NewFiberApp(settings APISettings, streamUseCase streamDelivery.UseCase, logger *slog.Logger) *FiberApp {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(fiberLogger.New())

	staticApp := fiber.New()
	app.Mount("/", staticApp)

	app.Get("/metrics", monitor.New(monitor.Config{Title: "bivi metrics page"}))
	app.Use(pprof.New())

	app.Get("/swagger/*", swagger.New(swagger.Config{
		DeepLinking:  false,
		DocExpansion: "none",
	}))

	staticApp.Use(streamDelivery.StaticHandler(streamUseCase, logger))
	staticApp.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))
	staticApp.Static("/", settings.ContentPath, fiber.Static{
		Compress:      true,
		CacheDuration: -1 * time.Second, // disable cache
	})

	api := app.Group("/api/v1")
	clientDelivery.NewDelivery(api.Group("/client"), settings.ClientLogPath, settings.UploadFilesizeLimit, logger)
	streamDelivery.NewDelivery(api.Group("/streams"), streamUseCase, logger)

	return &FiberApp{
		fiber:  app,
		logger: logger,
	}
}

func (f *FiberApp) Start(port string) error {
	return errors.Wrap(f.fiber.Listen(":"+port), "start application")
}

func (f *FiberApp) Shutdown(ctx context.Context) error {
	return errors.Wrap(f.fiber.ShutdownWithContext(ctx), "stop application")
}
