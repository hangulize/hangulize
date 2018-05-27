package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sublee/hangulize2/hangulize"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Hangulize 2",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("hangulize2-%s\n", hangulize.Version)
	},
}
