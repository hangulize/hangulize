package hre

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetaPatterns(t *testing.T) {
	assert.True(t, reLookbehind.MatchString(""))
	assert.True(t, reLookahead.MatchString(""))
}

func TestMacro(t *testing.T) {
	var p *Pattern

	p = fixturePattern(`@`) // @ means (a|e|i|o|u)
	assertFirstMatch(t, p, []string{
		o, "a",
		o, "ee",
		o, "iii",
		o, "no",
		o, "you",
		x, "sns", // no any vowel
	})

	p = fixturePattern(`_@_`)
	assertFirstMatch(t, p, []string{
		o, "_a_",
		x, "a__",
	})
}

func TestVars(t *testing.T) {
	var p *Pattern

	p = fixturePattern(`<abc>`)
	assertFirstMatch(t, p, []string{
		o, "a",
		o, "b",
		o, "c",
		x, "d",
	})

	p = fixturePattern(`<abc><def>`)
	assertFirstMatch(t, p, []string{
		o, "af",
		o, "bd",
		x, "db",
		o, "fcf",
		"    ^^",
	})
}

func TestSimple(t *testing.T) {
	p := fixturePattern(`hello, world`)
	assertFirstMatch(t, p, []string{
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
	assertFirstMatch(t, p, []string{
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
}

func TestLookahead(t *testing.T) {
	var p *Pattern

	p = fixturePattern(`han{gul}`)
	assertFirstMatch(t, p, []string{
		o, "hangul",
		"   ^^^   ",
		o, "hangulize",
		"   ^^^      ",
		x, "han",
		x, "hang",
		x, "hanja",
	})
}

func TestNegativeLookbehind(t *testing.T) {
	var p *Pattern

	p = fixturePattern(`{~han}gul`)
	assertFirstMatch(t, p, []string{
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
}

func TestNegativeLookahead(t *testing.T) {
	var p *Pattern

	p = fixturePattern(`han{~gul}`)
	assertFirstMatch(t, p, []string{
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
}

func TestLookaroundAndEdge(t *testing.T) {
	var p *Pattern

	p = fixturePattern(`foo{bar}$`)
	assertFirstMatch(t, p, []string{
		x, "foobar",
		x, " foobar ",
		x, "foo",
	})

	p = fixturePattern(`foo{~bar}$`)
	assertFirstMatch(t, p, []string{
		x, "foobar",
		x, " foobar ",
		o, "foo",
		"   ^^^",
	})

	p = fixturePattern(`^{foo}bar`)
	assertFirstMatch(t, p, []string{
		x, "foobar",
		x, " foobar ",
		x, "bar",
	})

	p = fixturePattern(`^{~foo}bar`)
	assertFirstMatch(t, p, []string{
		x, "foobar",
		x, " foobar ",
		o, "bar",
		"   ^^^",
	})
}

func TestLookaround(t *testing.T) {
	p := fixturePattern(`{ha}ng{ul}`)
	assertFirstMatch(t, p, []string{
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
	assertFirstMatch(t, p, []string{
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
	assertFirstMatch(t, p, []string{
		o, "foo",
		"   ^^^",
		o, "foobar",
		"   ^^^   ",
		o, "bar foobar",
		"       ^^^   ",
		x, "barfoobar",
	})

	p = fixturePattern(`^^foo`)
	assertFirstMatch(t, p, []string{
		o, "foobar",
		"   ^^^   ",
		o, "foobar bar",
		"   ^^^       ",
		x, "bar foobar",
	})

	p = fixturePattern(`foo$`)
	assertFirstMatch(t, p, []string{
		o, "foo",
		"   ^^^",
		o, "barfoo",
		"      ^^^",
		o, "barfoo foo",
		"      ^^^    ",
		x, "barfoobar",
	})

	p = fixturePattern(`foo$$`)
	assertFirstMatch(t, p, []string{
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
	assertFirstMatch(t, p, []string{
		o, "foo",
		"   ^^^",
		o, "foobar",
		"   ^^^   ",
		o, "bar foobar",
		"       ^^^   ",
		x, "barfoobar",
	})

	p = fixturePattern(`{^^}foo`)
	assertFirstMatch(t, p, []string{
		o, "foobar",
		"   ^^^   ",
		o, "foobar bar",
		"   ^^^       ",
		x, "bar foobar",
	})

	p = fixturePattern(`foo{$}`)
	assertFirstMatch(t, p, []string{
		o, "foo",
		"   ^^^",
		o, "barfoo",
		"      ^^^",
		o, "barfoo foo",
		"      ^^^    ",
		x, "barfoobar",
	})

	p = fixturePattern(`foo{$$}`)
	assertFirstMatch(t, p, []string{
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
	assertFirstMatch(t, p, []string{
		o, "acxxx",
		o, "xxbcx",
		x, "xxacx",
	})

	p = fixturePattern(`{foo}o{bar}`)
	assertFirstMatch(t, p, []string{
		o, "fooobar",
		"      ^   ",
	})
}

func TestMultipleNegativeLookahead(t *testing.T) {
	p := fixturePattern("foo{~foo}")
	w := "foobarbarfoofoo"
	assert.Equal(t, [][]int{[]int{0, 3}, []int{12, 15}}, p.Find(w, -1))
}

func TestMultipleNegativeLookbehind(t *testing.T) {
	p := fixturePattern("{~foo}foo")
	w := "barfoofoobarfoo"
	assert.Equal(t, [][]int{[]int{3, 6}, []int{12, 15}}, p.Find(w, -1))
}

func TestNegativeLookaroundWidth(t *testing.T) {
	p := fixturePattern("{~<vowels>}foo{~@}")
	naw, nbw := p.NegativeLookaroundWidths()
	assert.Equal(t, 1, naw)
	assert.Equal(t, 1, nbw)

	p = fixturePattern("{~<vowels>+}foo{~@*}")
	naw, nbw = p.NegativeLookaroundWidths()
	assert.Equal(t, -1, naw)
	assert.Equal(t, -1, nbw)
}

func TestMalformedPattern(t *testing.T) {
	p, err := NewPattern(`{a} {b} {c}`, nil, nil)
	assert.Error(t, err, p.Explain())
}

func TestBugs(t *testing.T) {
	var p *Pattern

	p = fixturePattern(`;|-`)
	assertFirstMatch(t, p, []string{
		o, "калинин,град-",
	})

	p = fixturePattern(`n{@|J}`)
	assertFirstMatch(t, p, []string{
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

func TestReplace(t *testing.T) {
	p := fixturePattern("{a}bc")
	rp := NewRPattern("_", nil, nil)

	assert.Equal(t, "a_bca_", p.Replace("abcbcabc", rp, -1))
}

func TestReplaceWithVars(t *testing.T) {
	vars := map[string][]string{
		"abc": []string{"a", "b", "c"},
		"xyz": []string{"x", "y", "z"},
	}

	p, _ := NewPattern("<abc>", nil, vars)
	rp := NewRPattern("<xyz>", nil, vars)

	assert.Equal(t, "xysolutely", p.Replace("absolutely", rp, -1))
}

func TestShiftedSubmatchIndex(t *testing.T) {
	p, _ := NewPattern("-|'", nil, nil)
	assert.Equal(t, [][]int{[]int{0, 1}}, p.Find("-", -1))
	assert.Equal(t, [][]int{[]int{0, 1}, []int{1, 2}}, p.Find("--", -1))
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

func BenchmarkLookaround(b *testing.B) {
	benchmark := func(title, expr string) {
		p := fixturePattern(expr)

		for i := 0; i <= 10; i++ {
			n := int(math.Pow(2, float64(i)))
			t := strings.Repeat("foobazfoofoobazbarbarbaz", n)

			b.Run(fmt.Sprintf("%s/%d[%d]", title, n, len(p.Find(t, -1))), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					p.Find(t, -1)
				}
			})
		}
	}

	benchmark("PA", "foo{foo}")  // O(n)
	benchmark("PB", "{foo}foo")  // O(n)
	benchmark("NA", "bar{~baz}") // O(n)
	benchmark("NB", "{~baz}bar") // O(n²)
}

// -----------------------------------------------------------------------------
// Examples

func ExamplePattern_Letters() {
	p, _ := NewPattern("^hello{,}", nil, nil)
	fmt.Println(p.Letters())
	// Output: [, e h l o]
}

func ExamplePattern_Find() {
	p, _ := NewPattern("^he(l+o){,}", nil, nil)
	fmt.Println(p.Find("hello, helo, hellllo", -1))
	// Output: [[0 5 2 5] [7 11 9 11]]
}

func ExamplePattern_Replace() {
	p, _ := NewPattern("foo{~bar}", nil, nil)
	rp := NewRPattern("xxx", nil, nil)
	fmt.Println(p.Replace("foo foobar foobaz", rp, -1))
	// Output: xxx foobar xxxbaz
}
