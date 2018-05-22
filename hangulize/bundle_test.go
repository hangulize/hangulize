package hangulize

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBundledTestCases(t *testing.T) {
	for _, lang := range ListLangs() {
		spec, ok := LoadSpec(lang)

		assert.Truef(t, ok, `failed to load "%s" spec`, lang)

		h := NewHangulizer(spec)

		for _, testCase := range spec.Test {
			loanword := testCase.Left()
			expected := testCase.Right()[0]

			ch := make(chan Event, 1000)
			got := h.HangulizeTrace(loanword, ch)

			if got == expected {
				continue
			}

			// Trace result to understand the failure reason.
			f := bytes.NewBufferString("")
			hr := strings.Repeat("-", 30)

			// Render failure message.
			fmt.Fprintln(f, hr)
			fmt.Fprintf(f, `lang: "%s"`, lang)
			fmt.Fprintln(f)
			fmt.Fprintf(f, `word: %#v`, loanword)
			fmt.Fprintln(f)
			fmt.Fprintln(f, hr)
			for e := range ch {
				fmt.Fprintln(f, e.String())
			}
			fmt.Fprintln(f, hr)

			assert.Equal(t, expected, got, f.String())

		}
	}
}
