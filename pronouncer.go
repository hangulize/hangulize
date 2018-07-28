package hangulize

// Pronouncer is an interface to guess pronunciation from spelling based on
// lexical analysis.
//
// The lexical analysis may require large size of dictionary data. To keep
// Hangulize lightweight, pronouncers are implemented out of this package.
//
// For example, there is the pronouncer for Furigana of Japanese in a separate
// package.
//
//   import "github.com/hangulize/hangulize"
//   import "github.com/hangulize/hangulize/pronounce/furigana"
//
//   hangulize.UsePronouncer(&furigana.P)
//   fmt.Println(hangulize.Hangulize("jpn", "日本語"))
//
type Pronouncer interface {
	ID() string
	Pronounce(string) string
}

// pronouncers is the registry holding the imported pronouncers.
var pronouncers = make(map[string]Pronouncer)

// UsePronouncer keeps a pronouncer for ready to use.
func UsePronouncer(p Pronouncer) bool {
	id := p.ID()

	if _, ok := pronouncers[id]; ok {
		return false
	}

	pronouncers[id] = p
	return true
}

// UnusePronouncer discards an imported pronouncer.
func UnusePronouncer(id string) bool {
	_, ok := pronouncers[id]
	if ok {
		delete(pronouncers, id)
	}
	return ok
}

// GetPronouncer returns the imported pronouncer by the ID.
func GetPronouncer(id string) (Pronouncer, bool) {
	p, ok := pronouncers[id]
	return p, ok
}
