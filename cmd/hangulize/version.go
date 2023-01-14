package main

import (
	"github.com/spf13/cobra"
)

var version string

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of Hangulize",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("hangulize %s", version)
		cmd.Println()
	},
}
