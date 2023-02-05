package main

import (
	"os"

	"github.com/hangulize/hangulize"
	"github.com/spf13/cobra"
)

var testCover bool
var testCoverProfile string

func init() {
	testCmd.Flags().BoolVarP(
		&testCover, "cover", "", false,
		"Enable coverage analysis.",
	)
	testCmd.Flags().StringVarP(
		&testCoverProfile, "coverprofile", "", "",
		"Write a coverage profile to the file after all tests have passed.",
	)

	rootCmd.AddCommand(testCmd)
}

var testCmd = &cobra.Command{
	Use:   "test HSL [HSL...]",
	Short: "Test examples in HSL files",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var cover *cover

		if testCoverProfile != "" {
			testCover = true
		}
		if testCover {
			cover = newCover()
		}

		failedAtLeastOnce := false

		for _, name := range args {
			// Open the HSL spec file.
			file, err := os.Open(name)
			if err != nil {
				return err
			}

			// Parse the spec.
			spec, err := hangulize.ParseSpec(file)
			if err != nil {
				return err
			}

			// Remember the name.
			cover.Visit(name)

			var (
				word     string
				expected string
				result   string
				traces   []hangulize.Trace
			)
			h := hangulize.New(spec)

			// Run test.
			for _, exm := range spec.Test {
				word, expected = exm[0], exm[1]

				if testCover {
					h.Trace(func(t hangulize.Trace) {
						traces = append(traces, t)
					})

					result, err = h.Hangulize(word)
					if err != nil {
						cmd.PrintErr(err)
					}
					for _, tr := range traces {
						if tr.Rule != nil {
							cover.Cover(name, tr.Step, tr.Rule.ID)
						}
					}
				} else {
					result, err = h.Hangulize(word)
					if err != nil {
						cmd.PrintErr(err)
					}
				}

				if result == expected {
					continue
				}

				// Test failed.
				cmd.Printf("%s: ", name)
				cmd.Printf(`"%s" -> "%s"`, word, result)
				cmd.Printf(`, expected: "%s"`, expected)
				cmd.Println()
				failedAtLeastOnce = true
			}
		}

		// Exit with 1 if failed at least once.
		if failedAtLeastOnce {
			os.Exit(1)
		}

		// Save the coverage profile.
		if testCoverProfile != "" {
			file, err := os.OpenFile(
				testCoverProfile,
				os.O_WRONLY|os.O_CREATE,
				0644,
			)
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
			defer file.Close()

			if err := cover.WriteProfile(file); err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		}

		// Print the coverage.
		if testCover {
			coverage := cover.Coverage()
			cmd.Printf("coverage: %.1f%% of rules\n", coverage*100)
		}

		return nil
	},
}
