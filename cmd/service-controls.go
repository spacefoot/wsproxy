//go:build windows

package main

import (
	"github.com/spacefoot/wsproxy/internal/windows"
	"github.com/spf13/cobra"
)

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
	serviceCmd.AddCommand(serviceStopCmd)
	serviceCmd.AddCommand(serviceStartCmd)
}
