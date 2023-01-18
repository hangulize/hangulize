package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/phonemize/furigana"
	"github.com/hangulize/hangulize/phonemize/pinyin"
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
	hangulize.ImportPhonemizer(&furigana.P)
	hangulize.ImportPhonemizer(&pinyin.P)
}

var rootCmd = &cobra.Command{
	Use:   "hangulize LANG WORD",
	Short: "Hangulize tools",

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lang := args[0]

		spec, ok := hangulize.LoadSpec(lang)
		if !ok {
			cmd.Println("Lang not supported:", lang)
			os.Exit(1)
		}

		h := hangulize.NewHangulizer(spec)
		hangulizeStream(cmd, args, h)
	},
}

func hangulizeStream(cmd *cobra.Command, args []string, h *hangulize.Hangulizer) {
	ch := make(chan string)
	go readWords(ch, args)

	for {
		word := <-ch
		if word == "" {
			break
		}

		if verbose {
			transcribed, traces, err := h.HangulizeTrace(word)
			if err != nil {
				cmd.PrintErr(err)
			}
			traces.Render(cmd.OutOrStderr())
			fmt.Fprintln(cmd.OutOrStdout(), transcribed)
		} else {
			transcribed, err := h.Hangulize(word)
			if err != nil {
				cmd.PrintErr(err)
			}
			fmt.Fprintln(cmd.OutOrStdout(), transcribed)
		}
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
