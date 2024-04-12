package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gitlab.teamdev.huds.su/bivi/backend/internal/models"
	"gitlab.teamdev.huds.su/bivi/backend/internal/pkg/app/errors"
	"log/slog"
)

type delivery struct {
	useCase InfoUseCase
	logger  *slog.Logger
}

func NewDelivery(api fiber.Router, useCase InfoUseCase, logger *slog.Logger) {
	handler := &delivery{
		useCase: useCase,
		logger:  logger,
	}
	api.Get("/", handler.getStreamsInfo)
	api.Get("/qualities", handler.getQualities)
}

// getQualities godoc
//
//	@Summary	Get stream qualities.
//	@Tags		Stream API
//	@Accept		json
//	@Success	200	{object}	QualitiesDTO		"List of available stream qualities"
//	@Failure	500	{object}	errors.FiberError	"Internal Server Error"
//	@Router		/streams/qualities [get]
func (d *delivery) getQualities(ctx *fiber.Ctx) error {
	var qualities []models.StreamQuality

	err := viper.UnmarshalKey("streams.qualities", &qualities)
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.NewFiberError(err.Error()))
	}

	resp := QualitiesDTO{Qualities: make([]QualityDTO, 0, len(qualities))}
	for _, quality := range qualities {
		resp.Qualities = append(resp.Qualities, QualityDTO{
			Height:               quality.Height,
			PreferredPeakBitrate: quality.PreferredPeakBitRate,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(resp)
}

// getStreamsInfo godoc
//
//	@Summary	Get stream descriptions.
//	@Tags		Stream API
//	@Accept		json
//	@Success	200	{object}	StreamsInfo			"List of available streams"
//	@Failure	500	{object}	errors.FiberError	"Internal Server Error"
//	@Router		/streams [get]
func (d *delivery) getStreamsInfo(ctx *fiber.Ctx) error {
	streams, err := d.useCase.GetStreamsInfo()
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(errors.NewFiberError(err.Error()))
	}

	return ctx.Status(fiber.StatusOK).JSON(StreamsInfo{Streams: streams})
}
