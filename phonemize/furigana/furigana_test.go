package furigana_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/phonemize/furigana"
)

func TestPhonemizer(t *testing.T) {
	assert.Implements(t, (*hangulize.Phonemizer)(nil), &furigana.P)
	assert.Equal(t, "furigana", furigana.P.ID())
}

func mustPhonemize(t *testing.T, word string) string {
	result, err := furigana.P.Phonemize(word)
	require.NoError(t, err)
	return result
}

func TestKana(t *testing.T) {
	assert.Equal(t, "ナイ", mustPhonemize(t, "ない"))
	assert.Equal(t, "ゲーム", mustPhonemize(t, "ゲーム"))
}

func TestRepeatKana(t *testing.T) {
	assert.Equal(t, "バナナ", mustPhonemize(t, "バナヽ"))
	assert.Equal(t, "イスズジドーシャ", mustPhonemize(t, "いすゞ自動車"))
}

func TestPersonNames(t *testing.T) {
	assert.Equal(t, "シンカイ マコト", mustPhonemize(t, "新海誠"))
	assert.Equal(t, "マツモト ユキヒロ", mustPhonemize(t, "松本行弘"))
}

func TestDeduplicateSpaces(t *testing.T) {
	assert.Equal(t, "シンカイ マコト", mustPhonemize(t, "新海 誠"))
	assert.Equal(t, "マツモト ユキヒロ", mustPhonemize(t, "松本 行弘"))
}

func TestProperNouns(t *testing.T) {
	assert.Equal(t, "トーイホッカイドー", mustPhonemize(t, "遠い北海道"))
}

func TestPunct(t *testing.T) {
	assert.Equal(t, "ヤサシイマリオ", mustPhonemize(t, "優しいマリオ"))
	assert.Equal(t, "ヤサシイ、マリオ", mustPhonemize(t, "優しい、マリオ"))
}

func TestVowel(t *testing.T) {
	assert.Equal(t, "ハウルノタベモノ", mustPhonemize(t, "ハウルの食べ物"))
	assert.Equal(t, "ハウルノウゴクシロ", mustPhonemize(t, "ハウルの動く城"))
}

func TestAmbiguousLinkedShortVowels(t *testing.T) {
	assert.Equal(t, "カワイイ", mustPhonemize(t, "可愛い"))
	assert.Equal(t, "オモウ", mustPhonemize(t, "思う"))
	assert.Equal(t, "ヌウ", mustPhonemize(t, "縫う"))
	assert.Equal(t, "キイテ", mustPhonemize(t, "聞いて"))
}

func TestLongVowelAcrossMorphemes(t *testing.T) {
	assert.Equal(t, "ハナサナカロー", mustPhonemize(t, "話さなかろう"))
}
