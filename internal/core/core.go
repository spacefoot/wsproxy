package core

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/spacefoot/wsproxy/internal/serial"
	"github.com/spacefoot/wsproxy/internal/serializer"
	"github.com/spacefoot/wsproxy/internal/static"
	"github.com/spacefoot/wsproxy/internal/websocket"
)

type Core struct {
	debug bool

	clientReader     chan []byte
	clientWriter     chan []byte
	clientRegistered chan bool

	scaleReader  chan []byte
	scaleWriter  chan []byte
	serialStatus chan serial.Status

	serializer serializer.Serializer
	hub        *websocket.Hub
	serial     serial.ISerial
}

func NewCore(debug bool) *Core {
	c := &Core{
		debug:            debug,
		clientReader:     make(chan []byte),
		clientWriter:     make(chan []byte),
		clientRegistered: make(chan bool),

		scaleReader:  make(chan []byte),
		scaleWriter:  make(chan []byte),
		serialStatus: make(chan serial.Status),

		serializer: serializer.Courier5000{},
	}

	c.hub = websocket.NewHub(c.clientReader, c.clientWriter, c.clientRegistered)

	if c.debug {
		slog.Info("running in debug mode")
		c.serial = serial.NewMock(c.serialStatus)
	} else {
		c.serial = serial.NewSerial(c.scaleReader, c.scaleWriter, c.serialStatus)
	}

	return c
}

func (c *Core) Run() {
	http.Handle("/ws", c.hub)
	http.HandleFunc("GET /{$}", c.serveDebugPage)

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
		if c.debug {
			go c.writeClient(data)
		}
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

func (c *Core) serveDebugPage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("index").Parse(static.Index))
	if err := t.Execute(w, static.IndexData{Debug: c.debug}); err != nil {
		slog.Error("error while rendering debug page", "err", err)
	}
}

func Run() {
	NewCore(false).Run()
}

func RunDebug(debug bool) {
	NewCore(debug).Run()
}
