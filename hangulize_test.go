package hangulize_test

import (
	"fmt"
	"testing"

	"github.com/hangulize/hangulize"
	"github.com/stretchr/testify/assert"
)

// TestLang generates subtests for bundled lang specs.
func TestLang(t *testing.T) {
	for _, lang := range hangulize.ListLangs() {
		spec, err := hangulize.LoadSpec(lang)
		assert.NoErrorf(t, err, `failed to load "%s" spec`, lang)

		for _, exm := range spec.Test {
			word := exm[0]
			expected := exm[1]

			t.Run(lang+"/"+word, func(t *testing.T) {
				assertHangulize(t, spec, expected, word)
			})
		}
	}
}

// -----------------------------------------------------------------------------
// Basic cases

func TestHangulizerSpec(t *testing.T) {
	spec, _ := hangulize.LoadSpec("ita")
	h := hangulize.New(spec)
	assert.Equal(t, spec, h.Spec())
}

// -----------------------------------------------------------------------------
// Edge cases

// TestSlash tests "/" in input word. The original Hangulize removes the "/" so
// the result was "글로르이아" instead of "글로르/이아".
func TestSlash(t *testing.T) {
	assert.Equal(t, "글로르/이아", mustHangulize(t, "ita", "glor/ia"))
	assert.Equal(t, "글로르{}이아", mustHangulize(t, "ita", "glor{}ia"))
}

func TestComma(t *testing.T) {
	assertHangulize(t, loadSpec("ita"), "글로르,이아", "glor,ia")
	assertHangulize(t, loadSpec("ita"), "콤,오", "com,o")
}

func TestPunctInVar(t *testing.T) {
	assertHangulize(t, loadSpec("nld"), "빔%", "Wim%")
	assertHangulize(t, loadSpec("cym"), "귀,림", "Gwi,lym")
	assertHangulize(t, loadSpec("wlm"), "카드,고데이", "Cad,Godeu")
}

func TestQuote(t *testing.T) {
	assert.Equal(t, "글로리아", mustHangulize(t, "ita", "glor'ia"))
	assert.Equal(t, "코모", mustHangulize(t, "ita", "com'o"))
}

func TestSpecials(t *testing.T) {
	assert.Equal(t, "<글로리아>", mustHangulize(t, "ita", "<gloria>"))
}

func TestHyphen(t *testing.T) {
	spec := mustParseSpec(`
	transcribe:
		"x" -> "-ㄱㅅ"
		"e-" -> "ㅣ"
		"e" -> "ㅔ"
	`)
	assert.Equal(t, "엑스야!", mustHangulizeSpec(t, spec, "ex야!"))
}

func TestDifferentAges(t *testing.T) {
	spec := mustParseSpec(`
	rewrite:
		"x" -> "xx"

	transcribe:
		"xx" -> "-ㄱㅅ"
		"e" -> "ㅔ"
	`)
	assert.Equal(t, "엑스야!", mustHangulizeSpec(t, spec, "ex야!"))
}

func TestKeepAndCleanup(t *testing.T) {
	spec := mustParseSpec(`
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
	assert.Equal(t, "ㅋ윽그스!", mustHangulizeSpec(t, spec, "ㅋ𐌄𐌗!"))
}

func TestSpace(t *testing.T) {
	spec := mustParseSpec(`
	rewrite:
		"van " -> "van/"

	transcribe:
		"van"  -> "반"
		"gogh" -> "고흐"
	`)
	assert.Equal(t, "반고흐", mustHangulizeSpec(t, spec, "van gogh"))
}

func TestZeroWidthSpace(t *testing.T) {
	spec := mustParseSpec(`
	rewrite:
		"a b" -> "a{}b"
		"^b"  -> "v"

	transcribe:
		"a" -> "ㅇ"
		"b" -> "ㅂ"
		"v" -> "ㅍ"
		"c" -> "ㅊ"
	`)
	assert.Equal(t, "으프 츠", mustHangulizeSpec(t, spec, "a b c"))
}

func TestVarToVar(t *testing.T) {
	spec := mustParseSpec(`
	vars:
		"abc" = "a", "b", "c"
		"def" = "d", "e", "f"
		"ghi" = "g", "h", "i"

	rewrite:
		"<abc><abc>" -> "<def><ghi>"

	transcribe:
		"a" -> "a"
		"b" -> "b"
		"c" -> "c"
		"d" -> "d"
		"e" -> "e"
		"f" -> "f"
		"g" -> "g"
		"h" -> "h"
		"i" -> "i"
	`)
	assert.Equal(t, "dg", mustHangulizeSpec(t, spec, "aa"))
	assert.Equal(t, "ei", mustHangulizeSpec(t, spec, "bc"))
}

func TestUnknownLang(t *testing.T) {
	result, err := hangulize.Hangulize("unknown", "hello")
	assert.ErrorIs(t, err, hangulize.ErrSpecNotFound)
	assert.Equal(t, "hello", result)
}

type stubTranslit struct{}

func (p *stubTranslit) Scheme() string {
	return "stub"
}

func (p *stubTranslit) Load() error {
	return nil
}

func (p *stubTranslit) Transliterate(word string) (string, error) {
	return "stub", nil
}

func TestInstanceTranslit(t *testing.T) {
	spec := mustParseSpec(`
	lang:
		id       = "test"
		codes    = "xx", "xxx"
		translit = "stub"

	transcribe:
		"stub" -> "스텁"
	`)
	h := hangulize.New(spec)

	_, err := h.Hangulize("1234")
	assert.ErrorIs(t, err, hangulize.ErrTranslitNotImported)

	h.UseTranslit(&stubTranslit{})
	result, err := h.Hangulize("1234")
	assert.NoError(t, err)
	assert.Equal(t, "스텁", result)

	h.UnuseTranslit("stub")
	_, err = h.Hangulize("1234")
	assert.ErrorIs(t, err, hangulize.ErrTranslitNotImported)
}

// -----------------------------------------------------------------------------
// Examples

func Example() {
	// Person names from http://iceager.egloos.com/2610028
	catalin, _ := hangulize.Hangulize("ron", "Cătălin Moroşanu")
	fmt.Println(catalin)

	jerrel, _ := hangulize.Hangulize("nld", "Jerrel Venetiaan")
	fmt.Println(jerrel)

	vitor, _ := hangulize.Hangulize("por", "Vítor Constâncio")
	fmt.Println(vitor)
	// Output:
	// 커털린 모로샤누
	// 예럴 페네티안
	// 비토르 콘스탄시우
}

func ExampleHangulize_cappuccino() {
	cappuccino, _ := hangulize.Hangulize("ita", "Cappuccino")
	fmt.Println(cappuccino)
	// Output: 카푸치노
}

func ExampleHangulize_nietzsche() {
	nietzsche, _ := hangulize.Hangulize("deu", "Friedrich Wilhelm Nietzsche")
	fmt.Println(nietzsche)
	// Output: 프리드리히 빌헬름 니체
}

func ExampleHangulize_shinkaiMakoto() {
	// import "github.com/hangulize/hangulize/translit"
	// translit.Install()

	shinkai, _ := hangulize.Hangulize("jpn", "新海誠")
	fmt.Println(shinkai)
	// Output: 신카이 마코토
}

func ExampleNew() {
	spec, _ := hangulize.LoadSpec("nld")
	h := hangulize.New(spec)

	gogh, _ := h.Hangulize("Vincent van Gogh")
	fmt.Println(gogh)
	// Output: 빈센트 반고흐
}
