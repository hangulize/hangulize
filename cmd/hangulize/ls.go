package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/hangulize/hangulize"
)

func init() {
	rootCmd.AddCommand(lsCmd)
}

var lsCmd = &cobra.Command{
	Use:   "ls [HSL...]",
	Short: "List of bundled specs or given HSLs",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		template := "%-8s %-8s %-24s %-24s\n"

		cmd.Printf(template, "LANG", "STAGE", "ENG", "KOR")

		for _, spec := range listSpecs(args) {
			cmd.Printf(template,
				spec.Lang.ID,
				spec.Config.Stage,
				spec.Lang.English,
				spec.Lang.Korean,
			)
		}
	},
}

func listSpecs(args []string) []*hangulize.Spec {
	specs := make([]*hangulize.Spec, 0)

	if len(args) == 0 {
		// Load bundled specs.
		for _, lang := range hangulize.ListLangs() {
			spec, _ := hangulize.LoadSpec(lang)
			specs = append(specs, spec)
		}
	} else {
		// Load specs from HSL files.
		for _, name := range args {
			file, err := os.Open(name)
			if err != nil {
				continue
			}

			spec, err := hangulize.ParseSpec(file)
			if err != nil {
				continue
			}

			specs = append(specs, spec)
		}
	}

	return specs
}
