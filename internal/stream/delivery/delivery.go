package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gitlab.teamdev.huds.su/bivi/backend/internal/models"
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
//	@Summary		Get stream qualities.
//	@Description	get stream qualities
//	@Tags			Stream API
//	@Accept			json
//	@Success		200
//	@Failure		500	{object}	map[string]string
//	@Router			/streams/qualities [get]
func (d *delivery) getQualities(ctx *fiber.Ctx) error {
	var qualities []models.StreamQuality

	err := viper.UnmarshalKey("streams.qualities", &qualities)
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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
//	@Summary		Get stream descriptions.
//	@Description	get stream descriptions
//	@Tags			Stream API
//	@Accept			json
//	@Success		200
//	@Failure		500	{object}	map[string]string
//	@Router			/streams [get]
func (d *delivery) getStreamsInfo(ctx *fiber.Ctx) error {
	streams, err := d.useCase.GetStreamsInfo()
	if err != nil {
		d.logger.Error(err.Error())
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(StreamsInfo{Streams: streams})
}
