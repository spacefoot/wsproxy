package main

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wstty",
	Short: "TTY WebSocket proxy",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
