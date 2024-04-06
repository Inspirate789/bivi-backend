package delivery

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
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
//	@Summary		Upload client logs.
//	@Description	upload client logs
//	@Tags			Client API
//	@Param			logs	formData	file	true	"Body with log file"
//	@Accept			json
//	@Success		200
//	@Failure		413	{object}	map[string]string
//	@Failure		422	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/client/logs [post]
func (d *delivery) postLogs(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("logs")
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	if file.Size > d.uploadFilesizeLimit {
		msg := fmt.Sprintf("The size of the uploaded file exceeds the limit: %d > %d", file.Size, d.uploadFilesizeLimit)
		d.logger.Error(msg)

		return ctx.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{"error": msg})
	}

	destination := d.logPath + "/" + file.Filename
	if err = ctx.SaveFile(file, destination); err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}
