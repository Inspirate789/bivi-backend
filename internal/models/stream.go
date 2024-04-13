package models

type StreamQuality struct {
	Width  uint `mapstructure:"width"`
	Height uint `mapstructure:"height"`
	// See more: https://developer.apple.com/documentation/avfoundation/avplayeritem/1388541-preferredpeakbitrate
	PreferredPeakBitRate uint `mapstructure:"preferredPeakBitRate"`
}

// StreamDescription godoc
//
//	@Description	Stream description (name and content paths)
//
// swagger:model
type StreamDescription struct {
	// Stream name
	Name string `example:"San Francisco" json:"name,omitempty"`
	// URL path to stream preview file
	PreviewPath string `example:"/content/U2FuIEZyYW5jaXNjbw==/preview.png" json:"previewPath,omitempty"`
	// URL path to stream HLS playlist file
	PlaylistPath string `example:"/content/U2FuIEZyYW5jaXNjbw==/playlist.m3u8" json:"playlistPath,omitempty"`
}
