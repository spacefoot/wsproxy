//go:build windows

package windows

import (
	"os"

	"github.com/spacefoot/wsproxy/internal/core"
	"golang.org/x/sys/windows/svc"
)

func RunIfService() (bool, error) {
	inService, err := svc.IsWindowsService()
	if err != nil {
		return false, err
	}

	if !inService {
		return false, nil
	}

	svc.Run(Name, &Handler{})
	return true, nil
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
