package main

import (
	"os"

	"github.com/hangulize/hangulize"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(devCmd)
}

var devCmd = &cobra.Command{
	Use:   "dev HSL",
	Short: "Develop a Hangulize spec",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]

		// Open an HSL file.
		file, err := os.Open(filename)
		if err != nil {
			return err
		}

		// Parse the spec.
		spec, err := hangulize.ParseSpec(file)
		if err != nil {
			return err
		}

		// Test the spec.
		h := hangulize.New(spec)
		hangulizeStream(cmd, args, h)
		return nil
	},
}
