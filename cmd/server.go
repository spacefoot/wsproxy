package main

import (
	"log/slog"

	"github.com/spacefoot/wsproxy/internal/core"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	RunE: func(cmd *cobra.Command, args []string) error {
		slog.SetLogLoggerLevel(slog.LevelDebug)
		core.Run()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
