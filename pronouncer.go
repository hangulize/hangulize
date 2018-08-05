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

// usePronouncer keeps a pronouncer into the given registry.
func usePronouncer(p Pronouncer, pronouncers *map[string]Pronouncer) bool {
	id := p.ID()

	if _, ok := (*pronouncers)[id]; ok {
		return false
	}

	(*pronouncers)[id] = p
	return true
}

// unusePronouncer discards a pronouncer from the given registry.
func unusePronouncer(id string, pronouncers *map[string]Pronouncer) bool {
	_, ok := (*pronouncers)[id]
	if ok {
		delete(*pronouncers, id)
	}
	return ok
}

// getPronouncer discards a pronouncer from the given registry.
func getPronouncer(
	id string,
	pronouncers *map[string]Pronouncer,
) (Pronouncer, bool) {
	p, ok := (*pronouncers)[id]
	return p, ok
}

// globalPronouncers is the registry holding the imported globalPronouncers.
var globalPronouncers = make(map[string]Pronouncer)

// UsePronouncer keeps a pronouncer for ready to use globally.
func UsePronouncer(p Pronouncer) bool {
	return usePronouncer(p, &globalPronouncers)
}

// UnusePronouncer discards a global pronouncer.
func UnusePronouncer(id string) bool {
	return unusePronouncer(id, &globalPronouncers)
}

// GetPronouncer returns a global pronouncer by the ID.
func GetPronouncer(id string) (Pronouncer, bool) {
	return getPronouncer(id, &globalPronouncers)
}
