package modules

import (
	"github.com/jabrach/telegram-admin-bot/cli"
	"log"
	"strings"
	"sync"
	"time"
)

type picUpdater struct {
	sync.Mutex
	chatID string
	queues map[int64]int64
}

const maxPicTimeout = 30

var PicUpdater = picUpdater{queues: map[int64]int64{}}

func (p *picUpdater) Update(msg *cli.Message, wrapper *cli.Wrapper) {
	if msg.Data.Event == "download" {
		log.Println("Downloaded some pic, setting as chat photo")
		wrapper.Exec("chat_set_photo", p.chatID, msg.Data.Result)
		return
	}

	if !filter(msg, IsMessage, FromManagedGroup) {
		return
	}

	if msg.Group().NoImages[msg.Data.From.PeerID] {
		return
	}

	if strings.TrimSpace(msg.Data.Text) == "/set_pic" {
		p.Lock()
		defer p.Unlock()
		p.queues[msg.Data.From.PeerID] = time.Now().Unix()
		p.chatID = msg.Data.To.ID
		log.Printf("Awaiting pic from %s\n", msg.Data.From.PrintName)
		return
	}

	if WithMedia(msg) && msg.Data.Media.Type == "photo" {
		if strings.TrimSpace(msg.Data.Media.Caption) == "/set_pic" {
			p.chatID = msg.Data.To.ID
			wrapper.Exec("load_photo", msg.ID)
			return
		}

		if t, ok := p.queues[msg.Data.From.PeerID]; ok {
			delta := time.Now().Unix() - t
			log.Printf("Got pic from %s after %v seconds\n", msg.Data.From.PrintName, delta)

			if delta < maxPicTimeout {
				log.Printf("Downloading pic from %s\n", msg.Data.From.PrintName)
				wrapper.Exec("load_photo", msg.ID)
			}
		}
	}
}
