package cli

type handlerFunc func(*Message, CLI)

type Message struct {
	ID      string   `json:"id"`
	Text    string   `json:"text"`
	Event   string   `json:"event"`
	From    *Someone `json:"from"`
	To      *Someone `json:"to"`
	Peer    *Someone `json:"peer"`
	Media   *Media   `json:"media"`
	Updates []string `json:"updates"`
	JSON    string   `json:"-"`
	Result  string   `json:"result"`
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
