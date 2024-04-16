//go:build windows

package main

import (
	"log/slog"

	"github.com/spacefoot/wsproxy/internal/windows"
	"github.com/spf13/cobra"
)

var serviceInstallCmd = &cobra.Command{
	Use: "install",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := windows.Install(); err != nil {
			return err
		}

		slog.Info("Service installed")

		if err := windows.Start(); err != nil {
			return err
		}

		slog.Info("Service started")
		return nil
	},
}

func init() {
	serviceCmd.AddCommand(serviceInstallCmd)
}
