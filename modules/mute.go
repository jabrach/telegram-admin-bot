package modules

import (
	"github.com/jabrach/telegram-admin-bot/cli"
	"log"
)

func Mute(msg *cli.Message, wrapper *cli.Wrapper) {
	if !filter(msg, IsMessage, FromManagedGroup) {
		return
	}

	log.Println("pennis dicke")

	if msg.Group().Muted[msg.Data.From.PeerID] {
		log.Printf("Message from muted user %v, deleting", msg.Data.From.PrintName)
		wrapper.Exec("delete_msg", msg.ID)
	} else if msg.Group().NoImages[msg.Data.From.PeerID] && WithMedia(msg) {
		log.Printf("Message with media from %v, deleting", msg.Data.From.PrintName)
		wrapper.Exec("delete_msg", msg.ID)
	}
}
