package hangulize

import "fmt"

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
