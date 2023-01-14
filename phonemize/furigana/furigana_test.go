package furigana_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/phonemize/furigana"
)

func TestPhonemizer(t *testing.T) {
	assert.Implements(t, (*hangulize.Phonemizer)(nil), &furigana.P)
	assert.Equal(t, "furigana", furigana.P.ID())
}

func TestKana(t *testing.T) {
	assert.Equal(t, "ナイ", furigana.P.Phonemize("ない"))
	assert.Equal(t, "ゲーム", furigana.P.Phonemize("ゲーム"))
}

func TestRepeatKana(t *testing.T) {
	assert.Equal(t, "バナナ", furigana.P.Phonemize("バナヽ"))
	assert.Equal(t, "イスズジドーシャ", furigana.P.Phonemize("いすゞ自動車"))
}

func TestPersonNames(t *testing.T) {
	assert.Equal(t, "シンカイ マコト", furigana.P.Phonemize("新海誠"))
	assert.Equal(t, "マツモト ユキヒロ", furigana.P.Phonemize("松本行弘"))
}

func TestDeduplicateSpaces(t *testing.T) {
	assert.Equal(t, "シンカイ マコト", furigana.P.Phonemize("新海 誠"))
	assert.Equal(t, "マツモト ユキヒロ", furigana.P.Phonemize("松本 行弘"))
}

func TestProperNouns(t *testing.T) {
	assert.Equal(t, "トーイホッカイドー", furigana.P.Phonemize("遠い北海道"))
}

func TestPunct(t *testing.T) {
	assert.Equal(t, "ヤサシイマリオ", furigana.P.Phonemize("優しいマリオ"))
	assert.Equal(t, "ヤサシイ、マリオ", furigana.P.Phonemize("優しい、マリオ"))
}

func TestVowel(t *testing.T) {
	assert.Equal(t, "ハウルノタベモノ", furigana.P.Phonemize("ハウルの食べ物"))
	assert.Equal(t, "ハウルノウゴクシロ", furigana.P.Phonemize("ハウルの動く城"))
}

func TestAmbiguousLinkedShortVowels(t *testing.T) {
	assert.Equal(t, "カワイイ", furigana.P.Phonemize("可愛い"))
	assert.Equal(t, "オモウ", furigana.P.Phonemize("思う"))
	assert.Equal(t, "ヌウ", furigana.P.Phonemize("縫う"))
	assert.Equal(t, "キイテ", furigana.P.Phonemize("聞いて"))
}

func TestLongVowelAcrossMorphemes(t *testing.T) {
	assert.Equal(t, "ハナサナカロー", furigana.P.Phonemize("話さなかろう"))
}
