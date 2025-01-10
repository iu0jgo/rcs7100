package engine

import (
	"crypto/tls"

	"github.com/iu0jgo/gumble/gumble"
	"github.com/iu0jgo/gumble/gumblemalgo"
)

type RCS7100 struct {
	Config              *gumble.Config
	Client              *gumble.Client
	Address             string
	TLSConfig           tls.Config
	ConnectAttempts     uint
	Stream              *gumblemalgo.Stream
	ChannelName         string
	IsConnected         bool
	IsTransmitting      bool
	PlaybackAudioDevice string
	CaptureAudioDevice  string
}
