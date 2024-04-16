//go:build windows

package main

import (
	"log/slog"

	"github.com/spacefoot/wsproxy/internal/windows"
	"github.com/spf13/cobra"
)

var serviceRemoveCmd = &cobra.Command{
	Use: "remove",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := windows.Remove(); err != nil {
			return err
		}

		slog.Info("Service scheduled for removal (after next reboot)")
		return nil
	},
}

func init() {
	serviceCmd.AddCommand(serviceRemoveCmd)
}
