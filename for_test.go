package hangulize

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// -----------------------------------------------------------------------------
// create or panic

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

func mustNewPattern(expr string, spec *Spec) *Pattern {
	p, err := newPattern(expr, spec.Macros, spec.Vars)
	if err != nil {
		panic(err)
	}
	return p
}

func mustNewRPattern(expr string, spec *Spec) *RPattern {
	p := newRPattern(expr, spec.Macros, spec.Vars)
	return p
}

// -----------------------------------------------------------------------------
// match assertion

const o = "MUST_MATCH"
const x = ""

// assertFirstMatch is a helper to test a pattern with multiple examples:
//
//  p := fixturePattern(`foo`)
//  assertFirstMatch(t, p, []string{
//    o, "foo",
//    "   ^^^",
//    o, "foobar",
//    "   ^^^   ",
//    x, "bar",
//  })
//
func assertFirstMatch(t *testing.T, p *Pattern, scenario []string) {
	drawUnderline := func(start int, stop int) string {
		return strings.Repeat(" ", start) + strings.Repeat("^", stop-start)
	}

	for i := 0; i < len(scenario); {
		mustMatch := scenario[i] == o
		word := scenario[i+1]
		i += 2

		matches := p.Find(word, 1)
		ok := len(matches) != 0

		m := []int{0, 0}
		if ok {
			m = matches[0]
		}

		if !mustMatch {
			if !ok {
				continue
			}

			assert.Failf(t, "unexpectedly matched", ""+
				"expected: NOT MATCH\n"+
				"actual  : \"%s\"\n"+
				"           %s\n"+
				"%s",
				word,
				drawUnderline(m[0], m[1]),
				p.Explain())
			continue
		}

		// Must match.
		if !ok {
			assert.Failf(t, "must match but not matched",
				"must MATCH with %#v\n%s", word, p.Explain())
		}

		if i == len(scenario) {
			break
		}

		// Find underline (^^^) which indicates expected match position.
		underline := scenario[i]
		if underline == o || underline == x {
			continue
		}
		i++

		if len(underline) != len(word)+3 {
			panic("underline length must be len(word)+3")
		}

		if len(m) == 0 {
			// Skip underline test because not matched.
			continue
		}

		start := strings.Index(underline, "^") - 3
		stop := strings.LastIndex(underline, "^") + 1 - 3

		expected := safeSlice(word, start, stop)
		got := word[m[0]:m[1]]

		assert.Equalf(t, expected, got, ""+
			"expected: \"%s\"\n"+
			"           %s\n"+
			"actual  : \"%s\"\n"+
			"           %s\n"+
			"%s",
			word, underline[3:],
			word, drawUnderline(m[0], m[1]),
			p.Explain())
	}
}

// -----------------------------------------------------------------------------
// hangulize assertion

func assertHangulize(t *testing.T, spec *Spec, expected string, word string) {
	h := NewHangulizer(spec)

	if h.Hangulize(word) == expected {
		return
	}

	// Trace only when failed to fast passing for most cases.
	got, tr := h.HangulizeTrace(word)

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
	for _, t := range tr {
		fmt.Fprintln(f, t.String())
	}
	fmt.Fprintln(f, hr)

	assert.Equal(t, expected, got, f.String())
}
