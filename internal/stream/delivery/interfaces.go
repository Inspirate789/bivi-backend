package delivery

import "gitlab.teamdev.huds.su/bivi/backend/internal/models"

type InfoUseCase interface {
	GetStreamsInfo() ([]models.StreamDescription, error)
}

type Streamer interface {
}
