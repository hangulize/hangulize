package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/sublee/hangulize2/hangulize"
)

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List of bundled lang specs",
	Run: func(cmd *cobra.Command, args []string) {
		template := "%-8s %-8s %-20s %-20s\n"

		fmt.Printf(template, "LANG", "STAGE", "ENG", "KOR")

		for _, lang := range hangulize.ListLangs() {
			spec, _ := hangulize.LoadSpec(lang)
			fmt.Printf(template,
				lang,
				spec.Config.Stage,
				spec.Lang.English,
				spec.Lang.Korean,
			)
		}
	},
}
