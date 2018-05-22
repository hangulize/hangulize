package hangulize

import (
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
	for i := 0; i < len(scenario); i += 2 {
		mustMatch := scenario[i] == o
		text := scenario[i+1]

		matched := p.Match(text)
		if mustMatch {
			assert.NotEmptyf(t, matched, "%s must match with %#v", p, text)
		} else {
			assert.Emptyf(t, matched, "%s must not match with  %#v", p, text)
		}
	}
}

func TestMacro(t *testing.T) {
	p = compile("@")
	assert.Equal(t, "(a|e|i|o|u)", p.reExpr)

	p = compile("sub@subl.ee")
	assert.Equal(t, "sub(a|e|i|o|u)subl.ee", p.reExpr)
}

func TestVars(t *testing.T) {
	p = compile("<abc>")
	assert.Equal(t, "(a|b|c)", p.reExpr)

	p = compile("<abc><def>")
	assert.Equal(t, "(a|b|c)(d|e|f)", p.reExpr)
}
