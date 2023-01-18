package hangulize_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/phonemize/furigana"
	"github.com/hangulize/hangulize/phonemize/pinyin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	hangulize.ImportPhonemizer(&furigana.P)
	hangulize.ImportPhonemizer(&pinyin.P)
}

func loadSpec(lang string) *hangulize.Spec {
	spec, ok := hangulize.LoadSpec(lang)
	if !ok {
		panic("failed to laod spec")
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

	actual, err := h.Hangulize(word)
	assert.NoError(t, err)

	if actual == expected {
		return
	}

	// Trace only when failed to fast passing for most cases.
	got, traces, err := h.HangulizeTrace(word)
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
	traces.Render(f)
	fmt.Fprintln(f, hr)

	assert.Equal(t, expected, got, f.String())
}
