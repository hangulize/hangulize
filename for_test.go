package hangulize_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/pkg/tracefmt"
	"github.com/hangulize/hangulize/translit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	translit.Install()
}

func loadSpec(lang string) *hangulize.Spec {
	spec, err := hangulize.LoadSpec(lang)
	if err != nil {
		panic(err)
	}
	return spec
}

func mustParseSpec(hsl string) *hangulize.Spec {
	spec, err := hangulize.ParseSpec(strings.NewReader(hsl))
	if err != nil {
		panic(err)
	}
	return spec
}

func mustHangulize(t *testing.T, lang, word string) string {
	result, err := hangulize.Hangulize(lang, word)
	require.NoError(t, err)
	return result
}

func mustHangulizeSpec(t *testing.T, spec *hangulize.Spec, word string) string {
	h := hangulize.New(spec)
	result, err := h.Hangulize(word)
	require.NoError(t, err)
	return result
}

func assertHangulize(t *testing.T, spec *hangulize.Spec, expected string, word string) {
	h := hangulize.New(spec)
	translit.Install(h)

	actual, err := h.Hangulize(word)
	assert.NoError(t, err)

	if actual == expected {
		return
	}

	// Trace only when failed to fast passing for most cases.
	traces := make([]hangulize.Trace, 0)
	h.Trace(func(t hangulize.Trace) {
		traces = append(traces, t)
	})

	got, err := h.Hangulize(word)
	assert.NoError(t, err)

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
	tracefmt.FprintTraces(f, traces)
	fmt.Fprintln(f, hr)

	assert.Equal(t, expected, got, f.String())
}
