//go:build windows

package windows

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/svc/mgr"
)

const (
	Name       = "wsproxy"
	Desciption = "TTY WebSocket proxy"
)

func exePath() (string, error) {
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}

	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		return "", fmt.Errorf("%s is directory", p)
	}

	if filepath.Ext(p) == "" {
		p += ".exe"
		fi, err := os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			return "", fmt.Errorf("%s is directory", p)
		}
	}

	return "", err
}

// Install registers the service with the Windows service control manager.
func Install() error {
	exepath, err := exePath()
	if err != nil {
		return err
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(Name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", Name)
	}

	s, err = m.CreateService(Name, exepath, mgr.Config{DisplayName: Desciption, StartType: mgr.StartAutomatic})
	if err != nil {
		return err
	}
	defer s.Close()

	return nil
}

// Remove unregisters the service with the Windows service control manager.
func Remove() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(Name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", Name)
	}
	defer s.Close()

	err = s.Delete()
	if err != nil {
		return err
	}

	return nil
}
