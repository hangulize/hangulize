package translit

import (
	"github.com/hangulize/hangulize"
	"github.com/hangulize/hangulize/translit/cyrillic"
	"github.com/hangulize/hangulize/translit/furigana"
	"github.com/hangulize/hangulize/translit/pinyin"
)

// Translits returns the standard Translits.
func Translits() []hangulize.Translit {
	ts := []hangulize.Translit{furigana.T, pinyin.T}
	ts = append(ts, cyrillic.Ts...)
	return ts
}

// Install imports all of the standard Translits.
func Install(h ...hangulize.Hangulizer) bool {
	if len(h) > 1 {
		panic("usage: translit.Use([hangulizer])")
	}

	useTranslit := hangulize.UseTranslit
	unuseTranslit := hangulize.UnuseTranslit
	if len(h) == 1 {
		useTranslit = h[0].UseTranslit
		unuseTranslit = h[0].UnuseTranslit
	}

	ts := Translits()
	for i, t := range ts {
		if ok := useTranslit(t); !ok {
			for j := 0; j < i; j++ {
				unuseTranslit(ts[j].Scheme())
			}
			return false
		}
	}
	return true
}
