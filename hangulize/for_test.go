package hangulize

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fxtSpec *Spec

func init() {
	// fixture spec
	spec, err := ParseSpec(strings.NewReader(strings.TrimSpace(`
# ------------------------------------------------------------------------------

vars:
	vowels = "a", "e", "i", "o", "u"
	abc    = "a", "b", "c"
	def    = "d", "e", "f"

macros:
	"@" = "<vowels>"

# ------------------------------------------------------------------------------
	`)))
	if err != nil {
		panic(err)
	}
	fxtSpec = spec
}

func fixturePattern(expr string) *Pattern {
	p, err := NewPattern(expr, fxtSpec.Macros, fxtSpec.Vars)
	if err != nil {
		panic(err)
	}

	return p
}

func fixtureRPatterns(exprs ...string) []*RPattern {
	to := make([]*RPattern, len(exprs))

	for i, expr := range exprs {
		p, err := NewRPattern(expr, fxtSpec.Macros, fxtSpec.Vars)
		if err != nil {
			panic(err)
		}
		to[i] = p
	}

	return to
}

const o = "MUST_MATCH"
const x = ""

// assertMatch is a helper to test a pattern with multiple examples:
//
//  p := fixturePattern(`foo`)
//  assertMatch(t, p, []string{
//    o, "foo",
//    "   ^^^",
//    o, "foobar",
//    "   ^^^   ",
//    x, "bar",
//  })
//
func assertMatch(t *testing.T, p *Pattern, scenario []string) {
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
				ExplainPattern(p))
			continue
		}

		// Must match.
		if !ok {
			assert.Failf(t, "must match but not matched",
				"must MATCH with %#v\n%s", word, ExplainPattern(p))
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
			ExplainPattern(p))
	}
}
