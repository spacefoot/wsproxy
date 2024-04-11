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
	clientReader     chan []byte
	clientWriter     chan []byte
	clientRegistered chan bool

	scaleReader  chan []byte
	scaleWriter  chan []byte
	serialStatus chan serial.Status

	serializer serializer.Serializer
	hub        *websocket.Hub
	serial     *serial.Serial
}

func NewCore() *Core {
	c := &Core{
		clientReader:     make(chan []byte),
		clientWriter:     make(chan []byte),
		clientRegistered: make(chan bool),

		scaleReader:  make(chan []byte),
		scaleWriter:  make(chan []byte),
		serialStatus: make(chan serial.Status),

		serializer: serializer.Courier5000{},
	}

	c.hub = websocket.NewHub(c.clientReader, c.clientWriter, c.clientRegistered)
	c.serial = serial.NewSerial(c.scaleReader, c.scaleWriter, c.serialStatus)

	return c
}

func (c *Core) Run() {
	http.Handle("/ws", c.hub)
	http.Handle("/", http.FileServerFS(static.FS))

	go c.run()
	go c.hub.Run()
	go c.serial.Run()

	slog.Info("server started", "addr", "http://localhost:23193")
	http.ListenAndServe("localhost:23193", nil)
}

func (c *Core) run() {
	for {
		select {
		case msg := <-c.clientReader:
			slog.Debug("received message", "from", "client", "data", string(msg))
			c.readClient(msg)
		case msg := <-c.scaleReader:
			slog.Debug("received message", "from", "scale", "data", string(msg))
			data, err := c.serializer.Read(msg)
			if err != nil {
				slog.Error("error while reading", "err", err)
				continue
			}
			if data != nil {
				c.writeClient(data)
			}
		case status := <-c.serialStatus:
			slog.Debug("received status", "from", "scale", "status", status)
			c.writeClient(&serializer.Status{
				Open: status.Open,
			})
		case <-c.clientRegistered:
			go c.serial.RequestStatus()
		}
	}
}

func (c *Core) readClient(msg []byte) {
	data, err := serializer.UnmarshalJSON(msg)
	if err != nil {
		slog.Error("error while unmarshalling", "err", err)
		return
	}

	switch data.(type) {
	case *serializer.RequestStatus:
		// RequestStatus will write to the serialStatus channel,
		// but the loop is already blocked in the clientReader path, causing a deadlock.
		// Must start another goroutine to avoid this.
		go c.serial.RequestStatus()
	case *serializer.Weight:
		// TMP echo back the weight for debugging
		go c.writeClient(data)
	default:
		slog.Warn("unknown message", "data", string(msg))
	}
}

func (c *Core) writeClient(msg any) {
	data, err := serializer.MarshalJSON(msg)
	if err != nil {
		slog.Error("error while marshalling", "err", err)
		return
	}

	if data == nil {
		return
	}

	slog.Debug("sending message", "to", "client", "data", string(data))
	c.clientWriter <- data
}

func Run() {
	NewCore().Run()
}
