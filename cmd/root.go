package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hangulize",
	Short: "Hangulize tools",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
