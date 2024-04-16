//go:build windows

package main

import (
	"log/slog"
	"os"

	"github.com/spacefoot/wsproxy/internal/core"
	"github.com/spacefoot/wsproxy/internal/windows"
	"golang.org/x/sys/windows/svc"
)

func main() {
	inService, err := svc.IsWindowsService()
	if err != nil {
		slog.Error("failed to determine if we are running in service", "err", err)
		os.Exit(1)
	}

	if inService {
		svc.Run(windows.Name, &Handler{})
		return
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type Handler struct{}

func (m *Handler) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	changes <- svc.Status{State: svc.StartPending}
	go core.Run()
	changes <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}
				os.Exit(0)
			}
		}
	}
}
