package hangulize

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func loadSpec(lang string) *Spec {
	spec, ok := LoadSpec(lang)
	if !ok {
		panic("failed to laod spec")
	}
	return spec
}

func mustParseSpec(hgl string) *Spec {
	spec, err := ParseSpec(strings.NewReader(hgl))
	if err != nil {
		panic(err)
	}
	return spec
}

func assertHangulize(t *testing.T, spec *Spec, expected string, word string) {
	h := NewHangulizer(spec)

	if h.Hangulize(word) == expected {
		return
	}

	// Trace only when failed to fast passing for most cases.
	got, traces := h.HangulizeTrace(word)

	// Trace result to understand the failure reason.
	f := bytes.NewBufferString("")
	hr := strings.Repeat("-", 30)

	// Render failure message.
	fmt.Fprintln(f, hr)
	fmt.Fprintf(f, `lang: "%s"`, spec.Lang.ID)
	fmt.Fprintln(f)
	fmt.Fprintf(f, `word: %#v`, word)
	fmt.Fprintln(f)
	fmt.Fprintln(f, hr)
	traces.Render(f)
	fmt.Fprintln(f, hr)

	assert.Equal(t, expected, got, f.String())
}
