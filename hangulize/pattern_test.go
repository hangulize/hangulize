package hangulize

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var spec *Spec
var p *Pattern

func init() {
	var err error
	spec, err = ParseSpec(strings.NewReader(strings.TrimSpace(`
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
}

func compile(expr string) *Pattern {
	return CompilePattern(expr, spec)
}

const o = "MUST_MATCH"
const x = ""

// assertMatch is a helper to test a pattern with multiple examples:
//
//  p := compile("foo")
//  assertMatch(t, p, []string{
//    o, "foo",
//    o, "foobar",
//    x, "bar",
//  })
//
func assertMatch(t *testing.T, p *Pattern, scenario []string) {
	info := fmt.Sprintf("re: /%s/, neg: /%s/", p.re, p.neg)

	for i := 0; i < len(scenario); i += 2 {
		mustMatch := scenario[i] == o
		text := scenario[i+1]

		matched := p.Match(text)

		if mustMatch {
			assert.NotEmptyf(t, matched,
				"%s must match with %#v\n%s", p, text, info)
		} else {
			assert.Emptyf(t, matched,
				"%s must not match with %#v\n%s", p, text, info)
		}
	}
}

func TestMacro(t *testing.T) {
	p = compile("@") // @ means (a|e|i|o|u)
	assertMatch(t, p, []string{
		o, "a",
		o, "ee",
		o, "iii",
		o, "no",
		o, "you",
		x, "sns", // no any vowel
	})

	p = compile("_@_")
	assertMatch(t, p, []string{
		o, "_a_",
		x, "a__",
	})
}

func TestVars(t *testing.T) {
	p = compile("<abc>")
	assertMatch(t, p, []string{
		o, "a",
		o, "b",
		o, "c",
		x, "d",
	})

	p = compile("<abc><def>")
	assertMatch(t, p, []string{
		o, "af",
		o, "bd",
		x, "db",
		o, "fcf",
	})
}

func TestSimple(t *testing.T) {
	p = compile("hello, world")
	assertMatch(t, p, []string{
		o, "hello, world",
		o, "__hello, world__",
		x, "bye, world",
	})
}

func TestLookbehind(t *testing.T) {
	p = compile("{han}gul")
	assertMatch(t, p, []string{
		o, "hangul",
		o, "hangulize",
		o, "__hangul",
		x, "gul",
		x, "ngul",
		x, "mogul",
	})

	p = compile("^{han}gul")
	assertMatch(t, p, []string{
		o, "hangul",
		o, "hangul__",
		x, "__hangul",
	})
}

func TestLookahead(t *testing.T) {
	p = compile("han{gul}")
	assertMatch(t, p, []string{
		o, "hangul",
		o, "hangulize",
		x, "han",
		x, "hang",
		x, "hanja",
	})

	p = compile("han{gul}$")
	assertMatch(t, p, []string{
		o, "hangul",
		o, "__hangul",
		x, "hangul__",
	})
}

func TestNegativeLookahead(t *testing.T) {
	p = compile("han{~gul}")
	assertMatch(t, p, []string{
		x, "hangul",
		x, "hangulize",
		o, "han",
		o, "hang",
		o, "hanja",
	})

	p = compile("han{~gul}$")
	assertMatch(t, p, []string{
		o, "han",
		o, "hangu",
		o, "han_gul",
		x, "hangul",
		x, "__hangul",
		x, "hangul__",
	})
}
