package delivery

import (
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"strings"
)

const pathPartsCount = 4

func StaticHandler(streamNameDecoder StreamNameDecoder, logger *slog.Logger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		path := ctx.Path()
		pathParts := strings.SplitN(path, "/", pathPartsCount)[:pathPartsCount]

		pathPrefix, encodedStreamName, filepath := pathParts[0]+"/"+pathParts[1], pathParts[2], pathParts[3]
		if filepath == "" { // structure doesn't match `/EncodedStreamName/Filepath`
			msg := "the requested file is not in any stream"
			logger.Error(msg)

			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
		}

		streamName, err := streamNameDecoder.DecodeString(encodedStreamName)
		if err != nil {
			logger.Error(err.Error())

			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		ctx.Path(pathPrefix + "/" + string(streamName) + "/" + filepath)

		return ctx.Next()
	}
}
