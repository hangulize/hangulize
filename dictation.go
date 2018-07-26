package hangulize

// Dictator is an interface to guess pronunciation from spelling based on
// lexical analysis. It would be used for the Normalize step.
//
// Some languages such as English or Japanese require external dictionaries to
// guess pronunciation from spelling. A dictator analyzes a word, based on
// external dictionaries, to pick out the lexemes that holding the source
// spelling and the dictated pronunciation. But external dictionaries may have
// large size of dataset. Therefore Hangulize itself doesn't include them by
// default due to it's lightness. A user would have a responsibility to install
// and use a dictator for the languages using a dictator.
//
// There is the Furigana dictator for Japanese in a separate package.
//
//   import "github.com/hangulize/hangulize"
//   import "github.com/hangulize/hangulize/dictate/furigana"
//
//   jpn := hangulize.LoadSpecWithDictator("jpn", furigana.Dictator)
//   h := hangulize.NewHangulizer(jpn)
//
//   fmt.Println(h.Hangulize("日本語"))
//
type Dictator interface {
	ID() string
	Dictate(string) [][2]string
}

// dictators is the registry holding the imported dictators.
var dictators = make(map[string]Dictator)

// UseDictator keeps a dictator for ready to use.
func UseDictator(d Dictator) bool {
	id := d.ID()

	if _, ok := dictators[id]; ok {
		return false
	}

	dictators[id] = d
	return true
}

// UnuseDictator discards an imported dictator.
func UnuseDictator(id string) bool {
	_, ok := dictators[id]
	if ok {
		delete(dictators, id)
	}
	return ok
}

// GetDictator returns the imported dictator by the ID.
func GetDictator(id string) (Dictator, bool) {
	d, ok := dictators[id]
	return d, ok
}
