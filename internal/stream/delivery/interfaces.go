package delivery

import "gitlab.teamdev.huds.su/bivi/backend/internal/models"

type UseCase interface {
	GetStreamsInfo() ([]models.StreamDescription, error)
}
