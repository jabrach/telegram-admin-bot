package modules

import (
	"github.com/jabrach/telegram-admin-bot/cli-wrapper"
)

type filterFunc func(*cli.Message) bool

func IsMessage(msg *cli.Message) bool {
	return msg.Data.Event == "message"
}

func IsUpdate(msg *cli.Message) bool {
	return msg.Data.Event == "updates"
}

func IsTitleUpdate(msg *cli.Message) bool {
	for _, upd := range msg.Data.Updates {
		if upd == "title" {
			return true
		}
	}
	return false
}

func FromManagedGroup(msg *cli.Message) bool {
	return msg.Group() != nil
}

func WithMedia(msg *cli.Message) bool {
	return msg.Data.Media != nil
}

func filter(msg *cli.Message, checks ...filterFunc) bool {
	for _, check := range checks {
		if !check(msg) {
			return false
		}
	}
	return true
}
