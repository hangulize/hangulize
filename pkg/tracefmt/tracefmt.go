package tracefmt

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hangulize/hangulize"
	"github.com/mattn/go-runewidth"
)

func PrintTraces(traces []hangulize.Trace) {
	FprintTraces(os.Stdout, traces)
}

func FprintTraces(w io.Writer, traces []hangulize.Trace) {
	step := ""
	columns := 0

	for _, t := range traces {
		if step != t.Step {
			step = t.Step
			fmt.Fprintf(w, "[%s]", step)
			fmt.Fprintln(w)
		}

		word := t.Word
		width := runewidth.StringWidth(word)
		if columns < width {
			columns = width/8*8 + 8
		}
		fmt.Fprintf(w, "  %s", word)
		fmt.Fprint(w, strings.Repeat(" ", columns-width))

		if t.Rule != nil {
			fmt.Fprintf(w, " | %s", t.Rule)
		} else if t.Why != "" {
			fmt.Fprintf(w, " | (%s)", t.Why)
		}

		fmt.Fprintln(w)
	}
}
