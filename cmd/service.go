//go:build windows

package main

import (
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage Windows service",
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
