package delivery

import "gitlab.teamdev.huds.su/bivi/backend/internal/models"

type InfoUseCase interface {
	GetStreamsInfo() ([]models.StreamDescription, error)
}

type StreamNameDecoder interface {
	DecodeString(s string) ([]byte, error)
}
