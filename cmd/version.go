package main

import (
	"encoding/json"

	"github.com/spacefoot/wsproxy/internal/core"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	RunE: func(cmd *cobra.Command, args []string) error {
		version := core.GetVersion()

		if v, _ := cmd.Flags().GetBool("json"); v {
			data, err := json.MarshalIndent(version, "", "  ")
			if err != nil {
				return err
			}
			cmd.Println(string(data))
			return nil
		}

		version.Print()
		return nil
	},
}

func init() {
	versionCmd.Flags().Bool("json", false, "Print in JSON format")
	rootCmd.AddCommand(versionCmd)
}
