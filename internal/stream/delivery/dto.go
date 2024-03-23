package delivery

import "gitlab.teamdev.huds.su/bivi/backend/internal/models"

type QualityDTO struct {
	Height               uint `json:"height,omitempty"`
	PreferredPeakBitrate uint `json:"preferredPeakBitRate,omitempty"`
}

type QualitiesDTO struct {
	Qualities []QualityDTO `json:"qualities,omitempty"`
}

type StreamsInfo struct {
	Streams []models.StreamDescription `json:"streams,omitempty"`
}
