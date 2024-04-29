//go:build windows

package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/inconshreveable/mousetrap"
	"github.com/spacefoot/wsproxy/internal/windows"
)

func main() {
	// If running from explorer.exe
	if mousetrap.StartedByExplorer() {
		if err := windows.Prompt(); err != nil {
			fmt.Println(err)
			time.Sleep(5 * time.Second)
			os.Exit(1)
		}
		return
	}

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
