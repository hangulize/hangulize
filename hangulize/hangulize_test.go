package hangulize

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLang generates subtests for bundled lang specs.
func TestLang(t *testing.T) {
	for _, lang := range ListLangs() {
		spec, ok := LoadSpec(lang)

		assert.Truef(t, ok, `failed to load "%s" spec`, lang)

		h := NewHangulizer(spec)

		for _, testCase := range spec.Test {
			loanword := testCase.Left()
			expected := testCase.Right()[0]

			t.Run(lang+"/"+loanword, func(t *testing.T) {
				got, tr := h.HangulizeTrace(loanword)
				if got == expected {
					return
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
				for _, t := range tr {
					fmt.Fprintln(f, t.String())
				}
				fmt.Fprintln(f, hr)

				assert.Equal(t, expected, got, f.String())
			})
		}
	}
}

// -----------------------------------------------------------------------------
// Edge cases

func hangulize(spec *Spec, word string) string {
	h := NewHangulizer(spec)
	return h.Hangulize(word)
}

// TestSlash tests "/" in input word.  The original Hangulize removes the "/"
// so the result was "글로르이아" instead of "글로르/이아".
func TestSlash(t *testing.T) {
	assert.Equal(t, "글로르/이아", Hangulize("ita", "glor/ia"))
	assert.Equal(t, "글로르{}이아", Hangulize("ita", "glor{}ia"))
}

func TestSpecials(t *testing.T) {
	assert.Equal(t, "<글로리아>", Hangulize("ita", "<gloria>"))
}

func TestHyphen(t *testing.T) {
	spec := parseSpec(`
	config:
		markers = "-"

	transcribe:
		"x" -> "-ㄱㅅ"
		"e-" -> "ㅣ"
		"e" -> "ㅔ"
	`)
	assert.Equal(t, "엑스야!", hangulize(spec, "ex야!"))
}

func TestDifferentAges(t *testing.T) {
	spec := parseSpec(`
	rewrite:
		"x" -> "xx"

	transcribe:
		"xx" -> "-ㄱㅅ"
		"e" -> "ㅔ"
	`)
	assert.Equal(t, "엑스야!", hangulize(spec, "ex야!"))
}

func TestKeepAndCleanup(t *testing.T) {
	spec := parseSpec(`
	rewrite:
		"𐌗"  -> "𐌗𐌗"
		"𐌄𐌗" -> "𐌊-"

	transcribe:
		"𐌊" -> "-ㄱ"
		"𐌗" -> "ㄱㅅ"
	`)
	// ㅋ𐌄 𐌗 !
	// ----│---------------------- rewrite
	//     ├─┐        𐌗->𐌗𐌗
	// ㅋ𐌄 𐌄 𐌗 !
	//   └┬┘
	//   ┌┴┐          𐌄𐌗->𐌊-
	// ㅋ𐌊 - 𐌗 !
	// --│------------------------ transcribe
	//   ├─┐          𐌊->ㄱ
	// ㅋ- ㄱ- 𐌗 !
	//         ├─┐    𐌗->-ㄱㅅ
	// ㅋ- ㄱ- ㄱㅅ!
	// ------│-------------------- cleanup
	//       x
	// ㅋ- ㄱㄱㅅ!
	// --├─┘┌┘┌┘------------------ jamo
	//   │ ┌┘┌┘
	// ㅋ윽그스!
	assert.Equal(t, "ㅋ윽그스!", hangulize(spec, "ㅋ𐌄𐌗!"))
}

func TestSpace(t *testing.T) {
	spec := parseSpec(`
	rewrite:
		"van " -> "van/"

	transcribe:
		"van"  -> "반"
		"gogh" -> "고흐"
	`)
	assert.Equal(t, "반고흐", hangulize(spec, "van gogh"))
}

func TestZeroWidthSpace(t *testing.T) {
	spec := parseSpec(`
	rewrite:
		"a b" -> "a{}b"
		"^b"  -> "v"

	transcribe:
		"a" -> "ㅇ"
		"b" -> "ㅂ"
		"v" -> "ㅍ"
		"c" -> "ㅊ"
	`)
	assert.Equal(t, "으프 츠", hangulize(spec, "a b c"))
}

// -----------------------------------------------------------------------------
// Benchmarks

func BenchmarkGloria(b *testing.B) {
	spec, _ := LoadSpec("ita")
	h := NewHangulizer(spec)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h.Hangulize("GLORIA")
	}
}
