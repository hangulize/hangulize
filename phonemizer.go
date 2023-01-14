package hangulize

import (
	"github.com/hangulize/hangulize/phonemize/furigana"
	"github.com/hangulize/hangulize/phonemize/pinyin"
)

// Register the built-in phonemizers.
func init() {
	if ok := UsePhonemizer(&pinyin.P); !ok {
		panic("failed to use phonemizer: pinyin")
	}

	if ok := UsePhonemizer(&furigana.P); !ok {
		panic("failed to use phonemizer: furigana")
	}
}

// Phonemizer is an interface to guess phonograms from a spelling based on
// lexical analysis.
type Phonemizer interface {
	ID() string
	Phonemize(string) string
}

// usePhonemizer keeps a phonemizer into the given registry.
func usePhonemizer(p Phonemizer, phonemizers map[string]Phonemizer) bool {
	id := p.ID()

	if _, ok := phonemizers[id]; ok {
		return false
	}

	phonemizers[id] = p
	return true
}

// unusePhonemizer discards a phonemizer from the given registry.
func unusePhonemizer(id string, phonemizers map[string]Phonemizer) bool {
	_, ok := phonemizers[id]
	if ok {
		delete(phonemizers, id)
	}
	return ok
}

// getPhonemizer discards a phonemizer from the given registry.
func getPhonemizer(id string, phonemizers map[string]Phonemizer) (Phonemizer, bool) {
	p, ok := phonemizers[id]
	return p, ok
}

// phonemizerRegistry is the registry holding the imported phonemizers.
var phonemizerRegistry = make(map[string]Phonemizer)

// UsePhonemizer imports a phonemizer for ready to use.
func UsePhonemizer(p Phonemizer) bool {
	return usePhonemizer(p, phonemizerRegistry)
}

// UnusePhonemizer discards a phonemizer by the ID.
func UnusePhonemizer(id string) bool {
	return unusePhonemizer(id, phonemizerRegistry)
}

// GetPhonemizer returns a phonemizer by the ID.
func GetPhonemizer(id string) (Phonemizer, bool) {
	return getPhonemizer(id, phonemizerRegistry)
}
