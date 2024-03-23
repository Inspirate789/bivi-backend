package usecase

import (
	"gitlab.teamdev.huds.su/bivi/backend/internal/models"
	"log/slog"
)

type UseCase struct {
	logger *slog.Logger
}

func NewUseCase(logger *slog.Logger) *UseCase {
	return &UseCase{
		logger: logger,
	}
}

func (uc *UseCase) GetStreamsInfo() ([]models.StreamDescription, error) {
	return []models.StreamDescription{
		{
			Name:         "Stream 1",
			PreviewPath:  "/preview",
			PlaylistPath: "/playlist",
		},
	}, nil
}
