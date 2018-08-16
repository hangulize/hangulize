package hangulize

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fixturePattern(expr string) *Pattern {
	spec := mustParseSpec(`
	vars:
		vowels = "a", "e", "i", "o", "u"
		abc    = "a", "b", "c"
		def    = "d", "e", "f"

	macros:
		"@" = "<vowels>"
	`)
	return mustNewPattern(expr, spec)
}

func TestMetaPatterns(t *testing.T) {
	assert.True(t, reLookbehind.MatchString(""))
	assert.True(t, reLookahead.MatchString(""))
}

func TestMacro(t *testing.T) {
	var p *Pattern

	p = fixturePattern(`@`) // @ means (a|e|i|o|u)
	assertMatch(t, p, []string{
		o, "a",
		o, "ee",
		o, "iii",
		o, "no",
		o, "you",
		x, "sns", // no any vowel
	})

	p = fixturePattern(`_@_`)
	assertMatch(t, p, []string{
		o, "_a_",
		x, "a__",
	})
}

func TestVars(t *testing.T) {
	var p *Pattern

	p = fixturePattern(`<abc>`)
	assertMatch(t, p, []string{
		o, "a",
		o, "b",
		o, "c",
		x, "d",
	})

	p = fixturePattern(`<abc><def>`)
	assertMatch(t, p, []string{
		o, "af",
		o, "bd",
		x, "db",
		o, "fcf",
		"    ^^",
	})
}

func TestSimple(t *testing.T) {
	p := fixturePattern(`hello, world`)
	assertMatch(t, p, []string{
		o, "hello, world",
		"   ^^^^^^^^^^^^",
		o, "__hello, world__",
		"     ^^^^^^^^^^^^  ",
		x, "bye, world",
	})
}

func TestLookbehind(t *testing.T) {
	var p *Pattern

	p = fixturePattern(`{han}gul`)
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

	p = fixturePattern(`^{han}gul`)
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
	var p *Pattern

	p = fixturePattern(`han{gul}`)
	assertMatch(t, p, []string{
		o, "hangul",
		"   ^^^   ",
		o, "hangulize",
		"   ^^^      ",
		x, "han",
		x, "hang",
		x, "hanja",
	})

	p = fixturePattern(`han{gul}$`)
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
	var p *Pattern

	p = fixturePattern(`{~han}gul`)
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

	p = fixturePattern(`^{~han}gul`)
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
	var p *Pattern

	p = fixturePattern(`han{~gul}`)
	assertMatch(t, p, []string{
		x, "hangul",
		x, "hangulize",
		o, "han",
		"   ^^^",
		o, "hang",
		"   ^^^ ",
		o, "hanja",
		"   ^^^  ",
		o, "hanhangul",
		"   ^^^      ",
		o, "han_hangul",
		"   ^^^       ",
	})

	p = fixturePattern(`han{~gul}$`)
	assertMatch(t, p, []string{
		o, "han",
		o, "hangu",
		o, "han_gul",
		x, "hangul",
		x, "__hangul",
		x, "hangul__",
	})
}

func TestLookaround(t *testing.T) {
	p := fixturePattern(`{ha}ng{ul}`)
	assertMatch(t, p, []string{
		o, "hangul",
		o, "hangulize",
		x, "ng",
		x, "hang",
		x, "ngul",
		x, "angu",
	})
}

func TestNegativeLookaround(t *testing.T) {
	p := fixturePattern(`{~ha}ng{~ul}`)
	assertMatch(t, p, []string{
		x, "hangul",
		x, "hangulize",
		x, "hang__",
		x, "__ngul",
		o, "__ng__",
	})
}

func TestEdge(t *testing.T) {
	var p *Pattern

	p = fixturePattern(`^foo`)
	assertMatch(t, p, []string{
		o, "foo",
		"   ^^^",
		o, "foobar",
		"   ^^^   ",
		o, "bar foobar",
		"       ^^^   ",
		x, "barfoobar",
	})

	p = fixturePattern(`^^foo`)
	assertMatch(t, p, []string{
		o, "foobar",
		"   ^^^   ",
		o, "foobar bar",
		"   ^^^       ",
		x, "bar foobar",
	})

	p = fixturePattern(`foo$`)
	assertMatch(t, p, []string{
		o, "foo",
		"   ^^^",
		o, "barfoo",
		"      ^^^",
		o, "barfoo foo",
		"      ^^^    ",
		x, "barfoobar",
	})

	p = fixturePattern(`foo$$`)
	assertMatch(t, p, []string{
		o, "barfoo",
		"      ^^^",
		o, "foo barfoo",
		"          ^^^",
		x, "foo foobar",
	})
}

func TestEdgeInLookaround(t *testing.T) {
	var p *Pattern

	p = fixturePattern(`{^}foo`)
	assertMatch(t, p, []string{
		o, "foo",
		"   ^^^",
		o, "foobar",
		"   ^^^   ",
		o, "bar foobar",
		"       ^^^   ",
		x, "barfoobar",
	})

	p = fixturePattern(`{^^}foo`)
	assertMatch(t, p, []string{
		o, "foobar",
		"   ^^^   ",
		o, "foobar bar",
		"   ^^^       ",
		x, "bar foobar",
	})

	p = fixturePattern(`foo{$}`)
	assertMatch(t, p, []string{
		o, "foo",
		"   ^^^",
		o, "barfoo",
		"      ^^^",
		o, "barfoo foo",
		"      ^^^    ",
		x, "barfoobar",
	})

	p = fixturePattern(`foo{$$}`)
	assertMatch(t, p, []string{
		o, "barfoo",
		"      ^^^",
		o, "foo barfoo",
		"          ^^^",
		x, "foo foobar",
	})
}

func TestComplexLookaround(t *testing.T) {
	var p *Pattern

	p = fixturePattern(`{^^a|b}c`)
	assertMatch(t, p, []string{
		o, "acxxx",
		o, "xxbcx",
		x, "xxacx",
	})

	p = fixturePattern(`{foo}o{bar}`)
	assertMatch(t, p, []string{
		o, "fooobar",
		"      ^   ",
	})
}

func TestMalformedPattern(t *testing.T) {
	p, err := newPattern(`{a} {b} {c}`, nil, nil)
	assert.Error(t, err, p.Explain())
}

func TestBugs(t *testing.T) {
	var p *Pattern

	p = fixturePattern(`;|-`)
	assertMatch(t, p, []string{
		o, "калинин,град-",
	})

	p = fixturePattern(`n{@|J}`)
	assertMatch(t, p, []string{
		o, "inJazio",
		"    ^     ",
	})
}

func TestEmptyPatternPanic(t *testing.T) {
	assert.Panics(t, func() { fixturePattern("") })
}

func TestZeroWidthMatchPanic(t *testing.T) {
	var p *Pattern

	p = fixturePattern("^")
	assert.Panics(t, func() { p.Find("abc", -1) })

	p = fixturePattern("$")
	assert.Panics(t, func() { p.Find("abc", -1) })

	p = fixturePattern("{a}")
	assert.Panics(t, func() { p.Find("abc", -1) })

	p = fixturePattern("{a}{b}")
	assert.Panics(t, func() { p.Find("abc", -1) })
}

// -----------------------------------------------------------------------------
// Benchmarks

func BenchmarkFind(b *testing.B) {
	p := fixturePattern("foo")

	for n := 0; n < 5; n++ {
		t := strings.Repeat("foo", int(math.Pow10(n)))

		b.Run(fmt.Sprintf("10**%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				p.Find(t, -1)
			}
		})
	}
}

func BenchmarkFindLookaround(b *testing.B) {
	p := fixturePattern("{~foo}foo")
	t := "barfoofoobarfoo"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		assert.Equal(b, [][]int{
			[]int{3, 6},
			[]int{12, 15},
		}, p.Find(t, -1))
	}
}

// -----------------------------------------------------------------------------
// Examples

func ExamplePattern_Letters() {
	p, _ := newPattern("^hello{,}", nil, nil)
	fmt.Println(p.Letters())
	// Output: [, e h l o]
}

func ExamplePattern_Find() {
	p, _ := newPattern("^he(l+o){,}", nil, nil)
	fmt.Println(p.Find("hello, helo, hellllo", -1))
	// Output: [[0 5 2 5] [7 11 9 11]]
}
