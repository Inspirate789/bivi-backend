package delivery

import "gitlab.teamdev.huds.su/bivi/backend/internal/models"

type QualityDTO struct {
	Height               uint `example:"720"     json:"height,omitempty"`
	PreferredPeakBitrate uint `example:"4200000" json:"preferredPeakBitRate,omitempty"`
}

type QualitiesDTO struct {
	Qualities []QualityDTO `json:"qualities,omitempty"`
}

type StreamsInfo struct {
	Streams []models.StreamDescription `json:"streams,omitempty"`
}
