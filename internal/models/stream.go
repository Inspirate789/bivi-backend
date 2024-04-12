package models

type StreamQuality struct {
	Width  uint `mapstructure:"width"`
	Height uint `mapstructure:"height"`
	// See more: https://developer.apple.com/documentation/avfoundation/avplayeritem/1388541-preferredpeakbitrate
	PreferredPeakBitRate uint `mapstructure:"preferredPeakBitRate"`
}

type StreamDescription struct {
	Name         string `example:"San Francisco"                               json:"name,omitempty"`
	PreviewPath  string `example:"/content/U2FuIEZyYW5jaXNjbw==/preview.png"   json:"previewPath,omitempty"`
	PlaylistPath string `example:"/content/U2FuIEZyYW5jaXNjbw==/playlist.m3u8" json:"playlistPath,omitempty"`
}
