package usecase

import (
	"encoding/base64"
	"github.com/spf13/viper"
	"gitlab.teamdev.huds.su/bivi/backend/internal/models"
	"log/slog"
	"path/filepath"
)

type UseCase struct {
	contentRoute      string
	repository        Repository
	streamNameEncoder StreamNameEncoder
	logger            *slog.Logger
}

func NewUseCase(
	contentRoute string,
	repository Repository,
	streamNameEncoder StreamNameEncoder,
	logger *slog.Logger,
) *UseCase {
	return &UseCase{
		contentRoute:      contentRoute,
		repository:        repository,
		streamNameEncoder: streamNameEncoder,
		logger:            logger,
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
			PreviewPath:  uc.contentRoute + "/" + filepath.Join(encodedName, viper.GetString("streams.filenames.preview")),
			PlaylistPath: uc.contentRoute + "/" + filepath.Join(encodedName, viper.GetString("streams.filenames.playlist")),
		})
	}

	return descriptions, nil
}
