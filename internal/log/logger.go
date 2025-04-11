package log

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

var Logger *slog.Logger

func init() {
	if strings.HasPrefix(os.Args[0], os.TempDir()) && strings.Contains(os.Args[0], "go-build") {
		Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		return
	}

	exe, err := filepath.Abs(os.Args[0])
	if err != nil {
		slog.Error("Failed to get absolute path", "error", err)
		return
	}

	fmt.Print(filepath.Join(filepath.Dir(exe), "wsproxy.log"))
	file, err := os.OpenFile(filepath.Join(filepath.Dir(exe), "wsproxy.log"), os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Failed to open log file", "error", err)
		return
	}

	Logger = slog.New(slog.NewJSONHandler(file, nil))
}
