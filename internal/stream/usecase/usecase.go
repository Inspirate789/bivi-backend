package usecase

import (
	"encoding/base64"
	"github.com/spf13/viper"
	"gitlab.teamdev.huds.su/bivi/backend/internal/models"
	"log/slog"
	"path/filepath"
)

type UseCase struct {
	repository Repository
	logger     *slog.Logger
}

func NewUseCase(repository Repository, logger *slog.Logger) *UseCase {
	return &UseCase{
		repository: repository,
		logger:     logger,
	}
}

func (uc *UseCase) GetStreamsInfo() ([]models.StreamDescription, error) {
	names, err := uc.repository.GetStreamNames()
	if err != nil {
		return nil, err
	}

	descriptions := make([]models.StreamDescription, 0, len(names))

	for _, name := range names {
		encodedName := base64.StdEncoding.EncodeToString([]byte(name))
		descriptions = append(descriptions, models.StreamDescription{
			Name:         name,
			PreviewPath:  "/" + filepath.Join(encodedName, viper.GetString("streams.filenames.preview")),
			PlaylistPath: "/" + filepath.Join(encodedName, viper.GetString("streams.filenames.playlist")),
		})
	}

	return descriptions, nil
}
