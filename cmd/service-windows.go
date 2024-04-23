//go:build windows

package main

import (
	"log/slog"

	"github.com/spacefoot/wsproxy/internal/windows"
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage Windows service",
}

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

func generateControlCommand(name string, handler func() error) *cobra.Command {
	return &cobra.Command{
		Use: name,
		RunE: func(cmd *cobra.Command, args []string) error {
			return handler()
		},
	}
}

var (
	serviceStartCmd = generateControlCommand("start", windows.Start)
	serviceStopCmd  = generateControlCommand("stop", windows.Stop)
)

func init() {
	rootCmd.AddCommand(serviceCmd)
	serviceCmd.AddCommand(serviceInstallCmd)
	serviceCmd.AddCommand(serviceRemoveCmd)
	serviceCmd.AddCommand(serviceStopCmd)
	serviceCmd.AddCommand(serviceStartCmd)
}
