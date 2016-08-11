package modules

import (
	"github.com/jabrach/telegram-admin-bot/cli-wrapper"
	"github.com/jabrach/telegram-admin-bot/config"
	"log"
	"strings"
	"sync"
)

type nameGuard struct {
	sync.Mutex
	topics map[int64]string
}

var Topic = nameGuard{}

func (n *nameGuard) Set(msg *cli.Message, wrapper cli.CLI) {
	if !filter(msg, IsMessage, FromManagedGroup) {
		return
	}
	text := strings.TrimSpace(msg.Data.Text)
	if len(text) < 8 || text[0:7] != "/topic " {
		return
	}
	topic := text[7:]
	log.Printf("Topic change from %s: %s", msg.Data.From.PrintName, topic)
	n.saveTopic(msg.Group().ID, topic)
	n.setTopic(msg.Data.To.ID, msg.Group(), wrapper)
}

func (n *nameGuard) Guard(msg *cli.Message, wrapper cli.CLI) {
	if !filter(msg, FromManagedGroup, IsUpdate, IsTitleUpdate) {
		return
	}
	if msg.Data.Peer == nil {
		return
	}

	if msg.Data.Peer.Title != n.fullTopic(msg.Group()) {
		log.Printf("Unwarranted topic change by %v, restoring", msg.Data.Peer.PrintName)
		n.setTopic(msg.Data.Peer.ID, msg.Group(), wrapper)
		wrapper.Exec("msg", msg.Data.Peer.ID, "Юзайте /topic, ущербы")
	}
}

func (n *nameGuard) saveTopic(groupID int64, topic string) {
	n.Lock()
	defer n.Unlock()
	if len(n.topics) == 0 {
		n.topics = map[int64]string{}
	}
	n.topics[groupID] = topic
}

func (n *nameGuard) setTopic(chatID string, group *config.Group, wrapper cli.CLI) {
	wrapper.Exec("rename_chat", chatID, n.fullTopic(group))
}

func (n *nameGuard) fullTopic(group *config.Group) string {
	return (group.Name + ". " + n.topics[group.ID])
}
