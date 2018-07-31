package furigana

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hangulize/hangulize"
)

func TestPronouncer(t *testing.T) {
	assert.Implements(t, (*hangulize.Pronouncer)(nil), &P)
}

func TestKana(t *testing.T) {
	assert.Equal(t, "ない", P.Pronounce("ない"))
	assert.Equal(t, "ゲーム", P.Pronounce("ゲーム"))
}

func TestRepeatKana(t *testing.T) {
	assert.Equal(t, "バナナ", P.Pronounce("バナヽ"))
	assert.Equal(t, "イスズジドウシャ", P.Pronounce("いすゞ自動車"))
}

func TestPersonNames(t *testing.T) {
	assert.Equal(t, "シンカイ マコト", P.Pronounce("新海誠"))
	assert.Equal(t, "マツモト ユキヒロ", P.Pronounce("松本行弘"))
}

func TestDeduplicateSpaces(t *testing.T) {
	assert.Equal(t, "シンカイ マコト", P.Pronounce("新海 誠"))
	assert.Equal(t, "マツモト ユキヒロ", P.Pronounce("松本 行弘"))
}

func TestProperNouns(t *testing.T) {
	assert.Equal(t, "トオイ・ホッカイドウ", P.Pronounce("遠い北海道"))
}

func TestPunct(t *testing.T) {
	assert.Equal(t, "ヤサシイ・マリオ", P.Pronounce("優しいマリオ"))
	assert.Equal(t, "ヤサシイ、マリオ", P.Pronounce("優しい、マリオ"))
}

func TestVowel(t *testing.T) {
	assert.Equal(t, "ハウルノタベモノ", P.Pronounce("ハウルの食べ物"))
	assert.Equal(t, "ハウルノ・ウゴクシロ", P.Pronounce("ハウルの動く城"))
}
