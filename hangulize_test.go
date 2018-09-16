package hangulize

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestLang generates subtests for bundled lang specs.
func TestLang(t *testing.T) {
	for _, lang := range ListLangs() {
		spec, ok := LoadSpec(lang)

		assert.Truef(t, ok, `failed to load "%s" spec`, lang)

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
	spec, _ := LoadSpec("ita")
	h := NewHangulizer(spec)
	assert.Equal(t, spec, h.Spec())
}

// -----------------------------------------------------------------------------
// Edge cases

func hangulize(spec *Spec, word string) string {
	h := NewHangulizer(spec)
	return h.Hangulize(word)
}

// TestSlash tests "/" in input word. The original Hangulize removes the "/" so
// the result was "글로르이아" instead of "글로르/이아".
func TestSlash(t *testing.T) {
	assert.Equal(t, "글로르/이아", Hangulize("ita", "glor/ia"))
	assert.Equal(t, "글로르{}이아", Hangulize("ita", "glor{}ia"))
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
	assert.Equal(t, "글로리아", Hangulize("ita", "glor'ia"))
	assert.Equal(t, "코모", Hangulize("ita", "com'o"))
}

func TestSpecials(t *testing.T) {
	assert.Equal(t, "<글로리아>", Hangulize("ita", "<gloria>"))
}

func TestHyphen(t *testing.T) {
	spec := mustParseSpec(`
	transcribe:
		"x" -> "-ㄱㅅ"
		"e-" -> "ㅣ"
		"e" -> "ㅔ"
	`)
	assert.Equal(t, "엑스야!", hangulize(spec, "ex야!"))
}

func TestDifferentAges(t *testing.T) {
	spec := mustParseSpec(`
	rewrite:
		"x" -> "xx"

	transcribe:
		"xx" -> "-ㄱㅅ"
		"e" -> "ㅔ"
	`)
	assert.Equal(t, "엑스야!", hangulize(spec, "ex야!"))
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
	assert.Equal(t, "ㅋ윽그스!", hangulize(spec, "ㅋ𐌄𐌗!"))
}

func TestSpace(t *testing.T) {
	spec := mustParseSpec(`
	rewrite:
		"van " -> "van/"

	transcribe:
		"van"  -> "반"
		"gogh" -> "고흐"
	`)
	assert.Equal(t, "반고흐", hangulize(spec, "van gogh"))
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
	assert.Equal(t, "으프 츠", hangulize(spec, "a b c"))
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
	assert.Equal(t, "dg", hangulize(spec, "aa"))
	assert.Equal(t, "ei", hangulize(spec, "bc"))
}

func TestUnknownLang(t *testing.T) {
	assert.Equal(t, "hello", Hangulize("unknown", "hello"))
}

type stubFurigana struct{}

func (p *stubFurigana) ID() string {
	return "furigana"
}

func (p *stubFurigana) Phonemize(word string) string {
	return "スタブ"
}

func TestInstancePhonemizers(t *testing.T) {
	spec, _ := LoadSpec("jpn")
	h := NewHangulizer(spec)

	h.UsePhonemizer(&stubFurigana{})
	assert.Equal(t, "스타부", h.Hangulize("1234"))

	h.UnusePhonemizer("furigana")
	assert.Equal(t, "1234", h.Hangulize("1234"))
}

// -----------------------------------------------------------------------------
// Examples

func Example() {
	// Person names from http://iceager.egloos.com/2610028
	fmt.Println(Hangulize("ron", "Cătălin Moroşanu"))
	fmt.Println(Hangulize("nld", "Jerrel Venetiaan"))
	fmt.Println(Hangulize("por", "Vítor Constâncio"))
	// Output:
	// 커털린 모로샤누
	// 예럴 페네티안
	// 비토르 콘스탄시우
}

func ExampleHangulize_cappuccino() {
	fmt.Println(Hangulize("ita", "Cappuccino"))
	// Output: 카푸치노
}

func ExampleHangulize_nietzsche() {
	fmt.Println(Hangulize("deu", "Friedrich Wilhelm Nietzsche"))
	// Output: 프리드리히 빌헬름 니체
}

func ExampleHangulize_shinkaiMakoto() {
	// import "github.com/hangulize/hangulize/phonemize/furigana"
	// UsePhonemizer(&furigana.P)

	fmt.Println(Hangulize("jpn", "新海誠"))
	// Output: 신카이 마코토
}

func ExampleNewHangulizer() {
	spec, _ := LoadSpec("nld")
	h := NewHangulizer(spec)

	fmt.Println(h.Hangulize("Vincent van Gogh"))
	// Output: 빈센트 반고흐
}
