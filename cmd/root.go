package main

import (
	"log/slog"
	"os"

	"github.com/spacefoot/wsproxy/internal/core"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wstty",
	Short: "TTY WebSocket proxy",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			slog.SetLogLoggerLevel(slog.LevelDebug)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		core.Run()
	},
}

func init() {
	rootCmd.PersistentFlags().Bool("verbose", false, "Verbose output")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
