package hangulize

import (
	"fmt"
	"testing"
)

// Here're all supported languages.
func ExampleListLangs() {
	for _, lang := range ListLangs() {
		fmt.Println(lang)
	}
	// Output:
	// aze
	// bel
	// bul
	// cat
	// ces
	// chi
	// cym
	// deu
	// ell
	// epo
	// est
	// fin
	// grc
	// hbs
	// hun
	// isl
	// ita
	// jpn
	// jpn-ck
	// kat-1
	// kat-2
	// lat
	// lav
	// lit
	// mkd
	// nld
	// pol
	// por
	// por-br
	// ron
	// rus
	// slk
	// slv
	// spa
	// sqi
	// swe
	// tur
	// ukr
	// vie
	// wlm
}

// -----------------------------------------------------------------------------
// Japanese

func TestJpnIgnoresLatin(t *testing.T) {
	assertHangulize(t, loadSpec("jpn"), "abc아", "abcあ")
}

func TestJpnKatakanaLongVowel(t *testing.T) {
	// http://www.textfugu.com/season-3/learn-katakana/2-3/
	jpn := loadSpec("jpn")

	assertHangulize(t, jpn, "위", "ウィ")
	assertHangulize(t, jpn, "웨", "ウェ")
	assertHangulize(t, jpn, "워", "ウォ")

	assertHangulize(t, jpn, "바", "ヴァ")
	assertHangulize(t, jpn, "비", "ヴィ")
	assertHangulize(t, jpn, "브", "ヴ")
	assertHangulize(t, jpn, "베", "ヴェ")
	assertHangulize(t, jpn, "보", "ヴォ")

	assertHangulize(t, jpn, "셰", "シェ")
	assertHangulize(t, jpn, "제", "ジェ")

	assertHangulize(t, jpn, "파", "ファ")
	assertHangulize(t, jpn, "피", "フィ")
	assertHangulize(t, jpn, "페", "フェ")
	assertHangulize(t, jpn, "포", "フォ")

	assertHangulize(t, jpn, "바", "ブァ")
	assertHangulize(t, jpn, "비", "ブィ")
	assertHangulize(t, jpn, "베", "ブェ")
	assertHangulize(t, jpn, "보", "ブォ")

	assertHangulize(t, jpn, "디", "ディ")
	assertHangulize(t, jpn, "도", "ドゥ")
	assertHangulize(t, jpn, "제", "ヂェ")

	assertHangulize(t, jpn, "디티", "ティティ")
	assertHangulize(t, jpn, "도토", "トゥトゥ")
	assertHangulize(t, jpn, "제체", "チェチェ")
}
