package main

import (
	"log/slog"

	"github.com/spacefoot/wsproxy/internal/core"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		debug, _ := cmd.Flags().GetBool("debug")
		simulateSerial, _ := cmd.Flags().GetBool("simulate-serial")

		if simulateSerial && !debug {
			debug = true
			slog.Info("Simulated serial requires debug mode. Enabling debug mode")
		}

		core.NewCore(core.CoreParams{
			Addr:           addr,
			Debug:          debug,
			SimulateSerial: simulateSerial,
		}).Run()
	},
}

func init() {
	serverCmd.Flags().Bool("debug", false, "Enable debug mode")
	serverCmd.Flags().Bool("simulate-serial", false, "Enable the simulated serial")
	serverCmd.Flags().String("addr", core.DefaultAddr, "Address to listen on")
	rootCmd.AddCommand(serverCmd)
}
