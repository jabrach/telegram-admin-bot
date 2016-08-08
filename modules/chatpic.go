package modules

import (
	"github.com/sthetz/tetanus/cli-wrapper"
	"github.com/sthetz/tetanus/config"
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

func (p *picUpdater) Update(msg *cli.Message, wrapper cli.CLI) {
	log.Println(msg.JSON)

	if msg.Event == "download" {
		log.Println("Downloaded some pic, setting as chat photo")
		wrapper.Exec("chat_set_photo", p.chatID, msg.Result)
		return
	}

	if !filter(msg, IsMessage, FromManagedGroup) {
		return
	}

	if config.NoImages(msg.From.PeerID) {
		return
	}

	if strings.TrimSpace(msg.Text) == "/set_pic" {
		p.Lock()
		defer p.Unlock()
		p.queues[msg.From.PeerID] = time.Now().Unix()
		p.chatID = msg.To.ID
		log.Printf("Awaiting pic from %s\n", msg.From.PrintName)
		return
	}

	if WithMedia(msg) && msg.Media.Type == "photo" {
		if t, ok := p.queues[msg.From.PeerID]; ok {
			delta := time.Now().Unix() - t
			log.Printf("Got pic from %s after %v seconds\n", msg.From.PrintName, delta)

			if delta < maxPicTimeout {
				log.Printf("Downloading pic from %s\n", msg.From.PrintName)
				wrapper.Exec("load_photo", msg.ID)
			}
		}
	}
}
