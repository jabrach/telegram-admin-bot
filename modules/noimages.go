package modules

import (
	"github.com/sthetz/tetanus/cli-wrapper"
	"github.com/sthetz/tetanus/config"
	"log"
)

func NoImages(msg *cli.Message, wrapper cli.CLI) {
	if msg.Event != "message" || msg.Media == nil {
		return
	}
	if msg.To.PeerType != "chat" || msg.To.PeerID != config.GroupID() {
		return
	}
	if !config.NoImages(msg.From.PeerID) {
		return
	}
	log.Printf("Message with media from %v, deleting", msg.From.PrintName)
	wrapper.Exec("delete_msg", msg.ID)
}
