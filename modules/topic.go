package modules

import (
	"github.com/sthetz/tetanus/cli-wrapper"
	"github.com/sthetz/tetanus/config"
	"log"
	"strings"
	"sync"
)

type nameGuard struct {
	sync.Mutex
	topic string
}

var Topic = nameGuard{}

func (n *nameGuard) Set(msg *cli.Message, wrapper cli.CLI) {
	if !filter(msg, IsMessage, FromManagedGroup) {
		return
	}
	text := strings.TrimSpace(msg.Text)
	if len(text) < 8 || text[0:7] != "/topic " {
		return
	}
	topic := text[7:]
	log.Printf("Topic change from %s: %s", msg.From.PrintName, topic)
	n.saveTopic(topic)
	n.setTopic(msg.To.ID, wrapper)
}

func (n *nameGuard) Guard(msg *cli.Message, wrapper cli.CLI) {
	if !filter(msg, FromManagedGroup, IsUpdate, IsTitleUpdate) {
		return
	}
	if msg.Peer == nil {
		return
	}

	if msg.Peer.Title != n.fullTopic() {
		log.Printf("Unwarranted topic change by %v, restoring", msg.Peer.PrintName)
		n.setTopic(msg.Peer.ID, wrapper)
		wrapper.Exec("msg", msg.Peer.ID, "Юзайте /topic, ущербы")
	}
}

func (n *nameGuard) saveTopic(topic string) {
	n.Lock()
	defer n.Unlock()
	n.topic = topic
}

func (n *nameGuard) setTopic(chatID string, wrapper cli.CLI) {
	wrapper.Exec("rename_chat", chatID, n.fullTopic())
}

func (n *nameGuard) fullTopic() string {
	return (config.GroupName() + ". " + n.topic)
}
