package websocket

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Hub struct {
	read  chan<- []byte
	write <-chan []byte

	peers      map[*Peer]bool
	register   chan *Peer
	unregister chan *Peer
}

func NewHub(read chan<- []byte, write <-chan []byte) *Hub {
	return &Hub{
		read:  read,
		write: write,

		peers:      map[*Peer]bool{},
		register:   make(chan *Peer),
		unregister: make(chan *Peer),
	}
}

func (h *Hub) Register(peer *Peer) {
	h.register <- peer
}

func (h *Hub) Unregister(peer *Peer) {
	h.unregister <- peer
}

func (h *Hub) Run() {
	for {
		select {
		case msg := <-h.write:
			for peer := range h.peers {
				peer.Write(msg)
			}
		case peer := <-h.register:
			h.peers[peer] = true
			slog.Debug("peer registered", "peer", peer.conn.RemoteAddr())
		case peer := <-h.unregister:
			delete(h.peers, peer)
			slog.Debug("peer unregistered", "peer", peer.conn.RemoteAddr())
		}
	}
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("upgrade error", "err", err)
		return
	}
	defer conn.Close()
	slog.Info("client connected", "addr", conn.RemoteAddr())

	peer := &Peer{
		conn: conn,
		read: h.read,
	}

	h.Register(peer)
	defer h.Unregister(peer)

	peer.Reader()
}
