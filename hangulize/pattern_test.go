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
//    "   ^^^",
//    o, "foobar",
//    "   ^^^   ",
//    x, "bar",
//  })
//
func assertMatch(t *testing.T, p *Pattern, scenario []string) {
	info := fmt.Sprintf("re: /%s/, neg: /%s/", p.re, p.neg)

	for i := 0; i < len(scenario); {
		mustMatch := scenario[i] == o
		text := scenario[i+1]
		i += 2

		matched := p.Match(text)

		if !mustMatch {
			assert.Emptyf(t, matched,
				"%s must NOT MATCH with %#v\n%s", p, text, info)
			continue
		}

		// Must match.
		assert.NotEmptyf(t, matched,
			"%s must MATCH with %#v\n%s", p, text, info)

		if i == len(scenario) {
			break
		}

		// Find underline (^^^) which indicates expected match position.
		underline := scenario[i]
		if underline == o || underline == x {
			continue
		}
		i++

		if len(underline) != len(text)+3 {
			panic("underline length must be len(text)+3")
		}

		if len(matched) == 0 {
			// Skip underline test because not matched.
			continue
		}

		start := strings.Index(underline, "^") - 3
		stop := strings.LastIndex(underline, "^") + 1 - 3

		expected := text[start:stop]
		got := text[matched[0]:matched[1]]

		assert.Equalf(t, expected, got,
			"%s on %#v must MATCH with %#v but %#v matched\n%s",
			p, text, expected, got, info)
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
		"    ^^",
	})
}

func TestSimple(t *testing.T) {
	p = compile("hello, world")
	assertMatch(t, p, []string{
		o, "hello, world",
		"   ^^^^^^^^^^^^",
		o, "__hello, world__",
		"     ^^^^^^^^^^^^  ",
		x, "bye, world",
	})
}

func TestLookbehind(t *testing.T) {
	p = compile("{han}gul")
	assertMatch(t, p, []string{
		o, "hangul",
		"      ^^^",
		o, "hangulize",
		"      ^^^   ",
		o, "__hangul",
		"        ^^^",
		x, "gul",
		x, "ngul",
		x, "mogul",
	})

	p = compile("^{han}gul")
	assertMatch(t, p, []string{
		o, "hangul",
		"      ^^^",
		o, "hangul__",
		"      ^^^  ",
		x, "__hangul",
		x, "__hangul__",
	})
}

func TestLookahead(t *testing.T) {
	p = compile("han{gul}")
	assertMatch(t, p, []string{
		o, "hangul",
		"   ^^^   ",
		o, "hangulize",
		"   ^^^      ",
		x, "han",
		x, "hang",
		x, "hanja",
	})

	p = compile("han{gul}$")
	assertMatch(t, p, []string{
		o, "hangul",
		"   ^^^   ",
		o, "__hangul",
		"     ^^^   ",
		x, "hangul__",
		x, "__hangul__",
	})
}

func TestNegativeLookbehind(t *testing.T) {
	p = compile("{~han}gul")
	assertMatch(t, p, []string{
		x, "hangul",
		x, "hangulize",
		x, "__hangul",

		o, "gul",
		"   ^^^",
		o, "ngul",
		"    ^^^",
		o, "mogul",
		"     ^^^",
		o, "hangulgul",
		"         ^^^",
		o, "hangul_gul",
		"          ^^^",
	})

	p = compile("^{~han}gul")
	assertMatch(t, p, []string{
		o, "gul",
		o, "angul",
		o, "han_gul",
		x, "hangul",
		x, "hangul__",
		x, "__hangul",
	})
}

func TestNegativeLookahead(t *testing.T) {
	p = compile("han{~gul}")
	assertMatch(t, p, []string{
		x, "hangul",
		x, "hangulize",
		o, "han",
		"   ^^^",
		o, "hang",
		"   ^^^ ",
		o, "hanja",
		"   ^^^  ",
		// o, "hanhangul",
		// "   ^^^      ",
		o, "han_hangul",
		"   ^^^       ",
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
