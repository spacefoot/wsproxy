//go:build windows

package main

import (
	"log/slog"
	"os"

	"github.com/spacefoot/wsproxy/internal/windows"
)

func main() {
	inService, err := windows.RunIfService()
	if err != nil {
		slog.Error("failed to run service", "err", err)
		return
	}

	if !inService {
		if err := rootCmd.Execute(); err != nil {
			os.Exit(1)
		}
	}
}
