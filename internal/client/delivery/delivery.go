package delivery

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gitlab.teamdev.huds.su/bivi/backend/internal/pkg/app/errors"
	"log/slog"
)

type delivery struct {
	logPath             string
	uploadFilesizeLimit int64
	logger              *slog.Logger
}

func NewDelivery(api fiber.Router, logPath string, uploadFilesizeLimit int64, logger *slog.Logger) {
	handler := &delivery{
		logPath:             logPath,
		uploadFilesizeLimit: uploadFilesizeLimit,
		logger:              logger,
	}
	api.Post("/logs", handler.postLogs)
}

// postLogs godoc
//
//	@Summary	Upload client logs.
//	@Tags		Client API
//	@Param		logs	formData	file	true	"Body with log file"
//	@Accept		json
//	@Success	200	"Log file saved"
//	@Failure	413	{object}	errors.FiberError	"Request file size is more than 5Mb"
//	@Failure	422	{object}	errors.FiberError	"Cannot get Multipart form file from request"
//	@Failure	500	{object}	errors.FiberError	"Internal Server Error"
//	@Router		/client/logs [post]
func (d *delivery) postLogs(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("logs")
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(errors.NewFiberError(err.Error()))
	}

	if file.Size > d.uploadFilesizeLimit {
		msg := fmt.Sprintf("The size of the uploaded file exceeds the limit: %d > %d", file.Size, d.uploadFilesizeLimit)
		d.logger.Error(msg)

		return ctx.Status(fiber.StatusRequestEntityTooLarge).JSON(errors.NewFiberError(msg))
	}

	destination := d.logPath + "/" + file.Filename
	if err = ctx.SaveFile(file, destination); err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.NewFiberError(err.Error()))
	}

	return ctx.SendStatus(fiber.StatusOK)
}
