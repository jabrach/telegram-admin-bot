package cli

import (
	"github.com/jabrach/telegram-admin-bot/config"
)

type msgHandler func(*Message, *Wrapper)

type MessageData struct {
	ID      string   `json:"id"`
	Text    string   `json:"text"`
	Event   string   `json:"event"`
	From    *Someone `json:"from"`
	To      *Someone `json:"to"`
	Peer    *Someone `json:"peer"`
	Media   *Media   `json:"media"`
	Updates []string `json:"updates"`
	Result  string   `json:"result"`
}

type Message struct {
	ID    string
	JSON  string
	Data  MessageData
	group *config.Group
}

type Someone struct {
	ID        string `json:"id"`
	PeerID    int64  `json:"peer_id"`
	PeerType  string `json:"peer_type"`
	PrintName string `json:"print_name"`
	Username  string `json:"username"`
	Title     string `json:"title"`
}

type Media struct {
	Type    string `json:"type"`
	Caption string `json:"caption"`
}

type Self struct {
	ID       string `json:"id"`
	PeerID   int64  `json:"peer_id"`
	Username string `json:"username"`
}

func (m *Message) Group() *config.Group {
	if m.group == nil {
		chat := m.Data.To
		if chat == nil {
			chat = m.Data.Peer
		}
		if chat == nil {
			return nil
		}
		m.group = config.ManagedGroup(chat.PeerID)
	}
	return m.group
}
