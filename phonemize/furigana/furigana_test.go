package furigana

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hangulize/hangulize"
)

func TestPhonemizer(t *testing.T) {
	assert.Implements(t, (*hangulize.Phonemizer)(nil), &P)
}

func TestKana(t *testing.T) {
	assert.Equal(t, "ない", P.Phonemize("ない"))
	assert.Equal(t, "ゲーム", P.Phonemize("ゲーム"))
}

func TestRepeatKana(t *testing.T) {
	assert.Equal(t, "バナナ", P.Phonemize("バナヽ"))
	assert.Equal(t, "イスズジドウシャ", P.Phonemize("いすゞ自動車"))
}

func TestPersonNames(t *testing.T) {
	assert.Equal(t, "シンカイ マコト", P.Phonemize("新海誠"))
	assert.Equal(t, "マツモト ユキヒロ", P.Phonemize("松本行弘"))
}

func TestDeduplicateSpaces(t *testing.T) {
	assert.Equal(t, "シンカイ マコト", P.Phonemize("新海 誠"))
	assert.Equal(t, "マツモト ユキヒロ", P.Phonemize("松本 行弘"))
}

func TestProperNouns(t *testing.T) {
	assert.Equal(t, "トオイ・ホッカイドウ", P.Phonemize("遠い北海道"))
}

func TestPunct(t *testing.T) {
	assert.Equal(t, "ヤサシイ・マリオ", P.Phonemize("優しいマリオ"))
	assert.Equal(t, "ヤサシイ、マリオ", P.Phonemize("優しい、マリオ"))
}

func TestVowel(t *testing.T) {
	assert.Equal(t, "ハウルノタベモノ", P.Phonemize("ハウルの食べ物"))
	assert.Equal(t, "ハウルノ・ウゴクシロ", P.Phonemize("ハウルの動く城"))
}

func TestAmbiguousLinkedShortVowels(t *testing.T) {
	assert.Equal(t, "カワイ・イ", P.Phonemize("可愛い"))
}
