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
