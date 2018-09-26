package hangulize

// Phonemizer is an interface to guess phonograms from a spelling based on
// lexical analysis.
//
// The lexical analysis may require large size of dictionary data. To keep
// Hangulize lightweight, phonemizers are implemented out of this package.
//
// For example, there is a phonemizer for Furigana of Japanese in a separate
// package.
//
//   import "github.com/hangulize/hangulize"
//   import "github.com/hangulize/hangulize/phonemize/furigana"
//
//   hangulize.UsePhonemizer(&furigana.P)
//   fmt.Println(hangulize.Hangulize("jpn", "日本語"))
//
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
func getPhonemizer(
	id string,
	phonemizers map[string]Phonemizer,
) (Phonemizer, bool) {
	p, ok := phonemizers[id]
	return p, ok
}

// phonemizerRegistry is the registry holding the imported phonemizerRegistry.
var phonemizerRegistry = make(map[string]Phonemizer)

// UsePhonemizer keeps a phonemizer for ready to use globally.
func UsePhonemizer(p Phonemizer) bool {
	return usePhonemizer(p, phonemizerRegistry)
}

// UnusePhonemizer discards a global phonemizer.
func UnusePhonemizer(id string) bool {
	return unusePhonemizer(id, phonemizerRegistry)
}

// GetPhonemizer returns a global phonemizer by the ID.
func GetPhonemizer(id string) (Phonemizer, bool) {
	return getPhonemizer(id, phonemizerRegistry)
}
