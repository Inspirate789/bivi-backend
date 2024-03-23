package models

type StreamQuality struct {
	Width  uint `mapstructure:"width"`
	Height uint `mapstructure:"height"`
	// See more: https://developer.apple.com/documentation/avfoundation/avplayeritem/1388541-preferredpeakbitrate
	PreferredPeakBitRate uint `mapstructure:"preferredPeakBitRate"`
}

type StreamDescription struct {
	Name         string `json:"name,omitempty"`
	PreviewPath  string `json:"previewPath,omitempty"`
	PlaylistPath string `json:"playlistPath,omitempty"`
}
