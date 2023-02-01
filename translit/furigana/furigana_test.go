package furigana_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hangulize/hangulize/translit/furigana"
)

func mustTransliterate(t *testing.T, word string) string {
	result, err := furigana.T.Transliterate(word)
	require.NoError(t, err)
	return result
}

func TestKana(t *testing.T) {
	assert.Equal(t, "ナイ", mustTransliterate(t, "ない"))
	assert.Equal(t, "ゲーム", mustTransliterate(t, "ゲーム"))
}

func TestRepeatKana(t *testing.T) {
	assert.Equal(t, "バナナ", mustTransliterate(t, "バナヽ"))
	assert.Equal(t, "イスズジドーシャ", mustTransliterate(t, "いすゞ自動車"))
}

func TestPersonNames(t *testing.T) {
	assert.Equal(t, "シンカイ マコト", mustTransliterate(t, "新海誠"))
	assert.Equal(t, "マツモト ユキヒロ", mustTransliterate(t, "松本行弘"))
}

func TestDeduplicateSpaces(t *testing.T) {
	assert.Equal(t, "シンカイ マコト", mustTransliterate(t, "新海 誠"))
	assert.Equal(t, "マツモト ユキヒロ", mustTransliterate(t, "松本 行弘"))
}

func TestProperNouns(t *testing.T) {
	assert.Equal(t, "トーイホッカイドー", mustTransliterate(t, "遠い北海道"))
}

func TestPunct(t *testing.T) {
	assert.Equal(t, "ヤサシイマリオ", mustTransliterate(t, "優しいマリオ"))
	assert.Equal(t, "ヤサシイ、マリオ", mustTransliterate(t, "優しい、マリオ"))
}

func TestVowel(t *testing.T) {
	assert.Equal(t, "ハウルノタベモノ", mustTransliterate(t, "ハウルの食べ物"))
	assert.Equal(t, "ハウルノウゴクシロ", mustTransliterate(t, "ハウルの動く城"))
}

func TestAmbiguousLinkedShortVowels(t *testing.T) {
	assert.Equal(t, "カワイイ", mustTransliterate(t, "可愛い"))
	assert.Equal(t, "オモウ", mustTransliterate(t, "思う"))
	assert.Equal(t, "ヌウ", mustTransliterate(t, "縫う"))
	assert.Equal(t, "キイテ", mustTransliterate(t, "聞いて"))
}

func TestLongVowelAcrossMorphemes(t *testing.T) {
	assert.Equal(t, "ハナサナカロー", mustTransliterate(t, "話さなかろう"))
}
