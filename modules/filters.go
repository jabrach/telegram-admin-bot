package modules

import (
	"github.com/sthetz/tetanus/cli-wrapper"
	"github.com/sthetz/tetanus/config"
)

type filterFunc func(*cli.Message) bool

func IsMessage(msg *cli.Message) bool {
	return msg.Event == "message"
}

func IsUpdate(msg *cli.Message) bool {
	return msg.Event == "updates"
}

func IsTitleUpdate(msg *cli.Message) bool {
	for _, upd := range msg.Updates {
		if upd == "title" {
			return true
		}
	}
	return false
}

func FromManagedGroup(msg *cli.Message) bool {
	chat := msg.To
	if chat == nil {
		chat = msg.Peer
	}
	if chat == nil {
		return false
	}

	return chat.PeerType == "chat" && chat.PeerID == config.GroupID()
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
