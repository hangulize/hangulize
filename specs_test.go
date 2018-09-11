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

func TestJpnAdditionalSounds(t *testing.T) {
	// http://www.textfugu.com/season-3/learn-katakana/2-3/
	// http://www.guidetojapanese.org/learn/grammar/katakana
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
	assertHangulize(t, jpn, "두", "ドゥ")
	assertHangulize(t, jpn, "제", "ヂェ")

	assertHangulize(t, jpn, "디티", "ティティ")
	assertHangulize(t, jpn, "두투", "トゥトゥ")
	assertHangulize(t, jpn, "제체", "チェチェ")
}

// -----------------------------------------------------------------------------
// Chinese

func TestChiUan(t *testing.T) {
	// From namu.wiki: https://namu.wiki/w/외래어%20표기법/중국어#s-3
	// u, uan, un이 j, q, x 뒤에 온다면 '위', '위안', '윈'으로 표기하고(예: ju
	// 쥐, quan 취안, xun 쉰), 다른 자음 뒤에 온다면 '우', '완', '운'으로
	// 표기한다(예: bu 부, duan 돤, hun 훈).
	chi := loadSpec("chi")

	assertHangulize(t, chi, "쥐", "ju")
	assertHangulize(t, chi, "취안", "quan")
	assertHangulize(t, chi, "쉰", "xun")
	assertHangulize(t, chi, "부", "bu")
	assertHangulize(t, chi, "돤", "duan")
	assertHangulize(t, chi, "훈", "hun")
}

func TestChiI(t *testing.T) {
	// From namu.wiki: https://namu.wiki/w/외래어%20표기법/중국어#s-3
	// i는 c, ch, r, s, sh, z, zh 뒤에 올 때는 ㅡ로 표기하며(예: chi 츠, ri 르,
	// shi 스, zi 쯔), 다른 자음 뒤에 올 때는 ㅣ로 표기한다(예: bi 비, ji 지,
	// li 리, ti 티).
	chi := loadSpec("chi")

	assertHangulize(t, chi, "츠", "chi")
	assertHangulize(t, chi, "르", "ri")
	assertHangulize(t, chi, "스", "shi")
	assertHangulize(t, chi, "쯔", "zi")
	assertHangulize(t, chi, "비", "bi")
	assertHangulize(t, chi, "지", "ji")
	assertHangulize(t, chi, "리", "li")
	assertHangulize(t, chi, "티", "ti")
}

func TestChiV(t *testing.T) {
	chi := loadSpec("chi")

	assertHangulize(t, chi, "뤼", "Lv")
	assertHangulize(t, chi, "뤼", "Lü")
}

func TestChiCJKUnified(t *testing.T) {
	chi := loadSpec("chi")

	assertHangulize(t, chi, "리", "李")
	assertHangulize(t, chi, "러", "樂")
}
