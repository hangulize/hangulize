package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/pkg/tracefmt"
	"github.com/hangulize/hangulize/translit"
	"github.com/spf13/cobra"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var verbose bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

var rootCmd = &cobra.Command{
	Use:   "hangulize LANG WORD",
	Short: "Hangulize tools",

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lang := args[0]

		spec, err := hangulize.LoadSpec(lang)
		if err != nil {
			cmd.Println("Lang not supported:", lang)
			os.Exit(1)
		}

		h := hangulize.New(spec)
		translit.Install(h)
		hangulizeStream(cmd, args, h)
	},
}

func hangulizeStream(cmd *cobra.Command, args []string, h hangulize.Hangulizer) {
	ch := make(chan string)
	go readWords(ch, args)

	for {
		word := <-ch
		if word == "" {
			break
		}

		var (
			result string
			traces []hangulize.Trace
			err    error
		)

		if verbose {
			h.Trace(func(t hangulize.Trace) {
				traces = append(traces, t)
			})
		}

		result, err = h.Hangulize(word)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
			return
		}

		tracefmt.FprintTraces(cmd.OutOrStderr(), traces)
		fmt.Fprintln(cmd.OutOrStdout(), result)
	}
}

func readWords(ch chan<- string, args []string) {
	if len(args) == 1 {
		reader := bufio.NewReader(os.Stdin)
		for {
			word, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			word = strings.TrimSpace(word)
			ch <- word
		}
	} else {
		for _, word := range args[1:] {
			if word != "" {
				ch <- word
			}
		}
	}
	ch <- ""
}
