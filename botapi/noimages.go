package botapi

import (
	"fmt"
	"github.com/sthetz/tetanus/config"
	"gopkg.in/telegram-bot-api.v4"
)

func (b *botInstance) noImages(msg *tgbotapi.Message) {
	if msg.Photo == nil && msg.Sticker == nil {
		return
	}
	if !config.NoImages(msg.From.ID) {
		return
	}

	if msg.From.UserName != "" {
		b.replyTo(msg, fmt.Sprintf("@%s ебало на 0", msg.From.UserName), false)
	}
}
