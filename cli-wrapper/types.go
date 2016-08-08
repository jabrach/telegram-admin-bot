package cli

type handlerFunc func(*Message, CLI)

type Message struct {
	ID    string   `json:"id"`
	Text  string   `json:"text"`
	Event string   `json:"event"`
	From  *Someone `json:"from"`
	To    *Someone `json:"to"`
	Media *Media   `json:"media"`
	JSON  string   `json:"-"`
}

type Someone struct {
	PeerID    int    `json:"peer_id"`
	PeerType  string `json:"peer_type"`
	PrintName string `json:"print_name"`
}

type Media struct {
	Type string `json:"type"`
}
