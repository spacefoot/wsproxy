package core

import (
	"encoding/json"
	"log/slog"

	"github.com/spacefoot/wsproxy/internal/serializer"
)

func (c *Core) weightDebug(msg []byte) {
	var w serializer.Weight
	if err := json.Unmarshal(msg, &w); err != nil {
		slog.Error("error while unmarshalling", "err", err)
		return
	}

	data, err := w.ToJSON()
	if err != nil {
		slog.Error("error while marshalling", "err", err)
		return
	}

	c.clientWriter <- data
}
