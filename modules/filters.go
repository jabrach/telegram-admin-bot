package modules

import (
	"github.com/sthetz/tetanus/cli-wrapper"
	"github.com/sthetz/tetanus/config"
)

type filterFunc func(*cli.Message) bool

func IsMessage(msg *cli.Message) bool {
	return msg.Event == "message"
}

func FromManagedGroup(msg *cli.Message) bool {
	return msg.To.PeerType == "chat" && msg.To.PeerID == config.GroupID()
}

func WithMedia(msg *cli.Message) bool {
	return msg.Media != nil
}

func filter(msg *cli.Message, checks ...filterFunc) bool {
	for _, check := range checks {
		if !check(msg) {
			return false
		}
	}
	return true
}