package log

import (
	"log/slog"
	"os"
	"path/filepath"
)

var discardLogger = slog.New(slog.DiscardHandler)
var file *os.File
var Logger *slog.Logger = discardLogger

func Enable() {
	if file != nil {
		return
	}

	exe, err := filepath.Abs(os.Args[0])
	if err != nil {
		slog.Error("Failed to get absolute path", "error", err)
		return
	}

	file, err = os.OpenFile(filepath.Join(filepath.Dir(exe), "wsproxy.log"), os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed to open log file", "error", err)
		return
	}

	Logger = slog.New(slog.NewJSONHandler(file, nil))
}

func Disable() {
	Logger = discardLogger
	if file != nil {
		file.Close()
		file = nil
	}
}
