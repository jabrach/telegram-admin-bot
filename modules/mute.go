package modules

import (
	"github.com/jabrach/telegram-admin-bot/cli-wrapper"
	"log"
)

func Mute(msg *cli.Message, wrapper cli.CLI) {
	if !filter(msg, IsMessage, FromManagedGroup, WithMedia) {
		return
	}

	if msg.Group().Muted[msg.Data.From.PeerID] {
		log.Printf("Message from muted user %v, deleting", msg.Data.From.PrintName)
		wrapper.Exec("delete_msg", msg.ID)
	} else if msg.Group().NoImages[msg.Data.From.PeerID] {
		log.Printf("Message with media from %v, deleting", msg.Data.From.PrintName)
		wrapper.Exec("delete_msg", msg.ID)
	}
}
