package delivery

import "gitlab.teamdev.huds.su/bivi/backend/internal/models"

// QualityDTO godoc
//
//	@Description	Stream quality (video height and bitrate)
//
// swagger:model
type QualityDTO struct {
	// Video height (px)
	// min: 360
	// max: 1080
	// example: 720
	Height uint `example:"720" json:"height,omitempty"`
	// Maximum bitrate of video in this quality (bytes per second)
	// min: 1400000
	// max: 8400000
	// example: 4200000
	PreferredPeakBitrate uint `example:"4200000" json:"preferredPeakBitRate,omitempty"`
}

// QualitiesDTO godoc
//
//	@Description	List of stream qualities (video height and bitrate)
//
// swagger:model
type QualitiesDTO struct {
	// List of stream qualities
	// min items: 1
	// max items: 3
	Qualities []QualityDTO `json:"qualities,omitempty"`
}

// StreamsInfo godoc
//
//	@Description	Stream descriptions (names and content paths)
//
// swagger:model
type StreamsInfo struct {
	// List of stream descriptions
	// min items: 0
	Streams []models.StreamDescription `json:"streams,omitempty"`
}
