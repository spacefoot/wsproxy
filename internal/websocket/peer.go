package websocket

import (
	"log/slog"

	"github.com/gorilla/websocket"
)

type Peer struct {
	conn *websocket.Conn
	read chan<- []byte
}

func (p *Peer) Write(msg []byte) {
	if err := p.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		slog.Error("error while writing", "err", err)
	}
}

func (p *Peer) Reader() {
	for {
		_, msg, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("error while reading", "err", err)
			}
			break
		}

		p.read <- msg
	}
}
