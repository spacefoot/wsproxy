package core

import (
	"log/slog"
	"net/http"

	"github.com/spacefoot/wsproxy/internal/serial"
	"github.com/spacefoot/wsproxy/internal/serializer"
	"github.com/spacefoot/wsproxy/internal/static"
	"github.com/spacefoot/wsproxy/internal/websocket"
)

type Core struct {
	clientReader chan []byte
	clientWriter chan []byte

	scaleReader chan []byte
	scaleWriter chan []byte

	serializer serializer.Serializer
}

func NewCore() *Core {
	return &Core{
		clientReader: make(chan []byte),
		clientWriter: make(chan []byte),

		scaleReader: make(chan []byte),
		scaleWriter: make(chan []byte),

		serializer: serializer.Courier5000{},
	}
}

func (c *Core) Run() {
	hub := websocket.NewHub(c.clientReader, c.clientWriter)
	scale := serial.NewSerial(c.scaleReader, c.scaleWriter)

	http.Handle("/ws", hub)
	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(static.TestPage)
	})

	go c.run()
	go hub.Run()
	go scale.Run()

	slog.Info("server started", "addr", "http://localhost:23193")
	http.ListenAndServe("[::1]:23193", nil)
}

func (c *Core) run() {
	for {
		select {
		case msg := <-c.clientReader:
			slog.Debug("received message", "from", "client", "data", string(msg))
			c.weightDebug(msg) // TMP
		case msg := <-c.scaleReader:
			slog.Debug("received message", "from", "scale", "data", string(msg))
			data, err := c.serializer.Read(msg)
			if err != nil {
				slog.Error("error while reading", "err", err)
				continue
			}
			if data != nil {
				c.clientWriter <- data
			}
		}
	}
}

func Run() {
	NewCore().Run()
}
