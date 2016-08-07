package botapi

import (
	"fmt"
	"github.com/sthetz/tetanus/config"
	"gopkg.in/telegram-bot-api.v4"
)

const longPollingTimeout = 60

type BotAPI interface {
	Listen()
}

type botInstance struct {
	apiToken string
	api      *tgbotapi.BotAPI
}

func New(apiToken string) BotAPI {
	bot := &botInstance{}
	err := bot.InitAPI(apiToken)

	if err != nil {
		panic(err)
	}

	return bot
}

func (b *botInstance) InitAPI(apiToken string) error {
	b.apiToken = apiToken
	api, err := tgbotapi.NewBotAPI(b.apiToken)
	b.api = api

	return err
}

func (b *botInstance) Listen() {
	fmt.Printf("Authorized on account @%v\n", b.api.Self.UserName)

	updates, err := b.api.GetUpdatesChan(b.getUpdatesCfg())
	if err != nil {
		panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		b.processMessage(update.Message)
	}
}

func (b *botInstance) replyTo(msg *tgbotapi.Message, text string, citate bool) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	if citate {
		reply.ReplyToMessageID = msg.MessageID
	}
	b.api.Send(reply)
}

func (b *botInstance) processMessage(msg *tgbotapi.Message) {
	if msg.Chat.ID != config.GroupID() {
		b.leaveChat(msg)
		return
	}
	fmt.Printf("%s (%v): %s\n", msg.From.String(), msg.From.ID, msg.Text)

	b.noImages(msg)
}

func (b *botInstance) leaveChat(msg *tgbotapi.Message) {
	if msg.Chat != nil && msg.Chat.IsGroup() {
		fmt.Printf("Message from chat #%v, leaving", msg.Chat.ID)
		b.api.LeaveChat(msg.Chat.ChatConfig())
	}
}

func (b *botInstance) getUpdatesCfg() tgbotapi.UpdateConfig {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = longPollingTimeout
	return u
}
