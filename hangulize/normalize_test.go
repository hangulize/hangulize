package hangulize

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeRoman(t *testing.T) {
	roman := &RomanNormalizer{}

	assert.Equal(t, "hello", Normalize("Hello", roman, nil))
	assert.Equal(t, "cafe", Normalize("Café", roman, nil))
	assert.Equal(t, "melee", Normalize("Mêlée", roman, nil))
	assert.Equal(t, " cafe, hi! ", Normalize(" Café, Hi! ", roman, nil))
}

func TestNormalizeKana(t *testing.T) {
	kana := &KanaNormalizer{}

	assert.Equal(t, "ア", Normalize("あ", kana, nil))
	assert.Equal(t, "ァ", Normalize("ぁ", kana, nil))
}

// -----------------------------------------------------------------------------
// Examples

func ExampleNormalize() {
	roman := &RomanNormalizer{}
	fmt.Println(Normalize("Café", roman, nil))
	// Output: cafe
}

// Letters in keep will not normalized by the given normalizer.
func ExampleNormalize_keep() {
	roman := &RomanNormalizer{}
	fmt.Println(Normalize("pisáť", roman, []string{"ť"}))
	// Output: pisať
}

func ExampleRomanNormalizer() {
	roman := &RomanNormalizer{}
	fmt.Println(Normalize("Lūciī a fīliī", roman, nil))
	// Output: lucii a filii
}

func ExampleKanaNormalizer() {
	kana := &KanaNormalizer{}
	fmt.Println(Normalize("こんにちは", kana, nil))
	// Output: コンニチハ
}

func ExampleGetNormalizer() {
	kana, _ := GetNormalizer("kana")
	Normalize("こんにちは", kana, nil)
}
