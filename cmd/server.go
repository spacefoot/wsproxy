package main

import (
	"github.com/spacefoot/wsproxy/internal/core"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Flags().GetBool("debug")
		core.RunDebug(debug)
	},
}

func init() {
	serverCmd.Flags().Bool("debug", false, "Enable debug mode")
	rootCmd.AddCommand(serverCmd)
}
