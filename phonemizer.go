package hangulize

// Phonemizer is an interface to guess phonograms from a spelling based on
// lexical analysis.
type Phonemizer interface {
	ID() string
	Load() error
	Phonemize(string) (string, error)
}

// importPhonemizer keeps a phonemizer into the given registry.
func importPhonemizer(p Phonemizer, registry map[string]Phonemizer) bool {
	id := p.ID()

	if _, ok := registry[id]; ok {
		return false
	}

	registry[id] = p
	return true
}

// unimportPhonemizer discards a phonemizer from the given registry.
func unimportPhonemizer(id string, registry map[string]Phonemizer) bool {
	_, ok := registry[id]
	if ok {
		delete(registry, id)
	}
	return ok
}

// getPhonemizer finds a phonemizer from the given registry.
func getPhonemizer(id string, registry map[string]Phonemizer) (Phonemizer, bool) {
	p, ok := registry[id]
	return p, ok
}

// defaultPhonemizerRegistry is the registry holding the imported phonemizers.
var defaultPhonemizerRegistry = make(map[string]Phonemizer)

// ImportPhonemizer imports a phonemizer for ready to use.
func ImportPhonemizer(p Phonemizer) bool {
	return importPhonemizer(p, defaultPhonemizerRegistry)
}

// UnimportPhonemizer discards a phonemizer by the ID.
func UnimportPhonemizer(id string) bool {
	return unimportPhonemizer(id, defaultPhonemizerRegistry)
}

// DefaultPhonemizer returns a phonemizer from the default registry.
func DefaultPhonemizer(id string) (Phonemizer, bool) {
	return getPhonemizer(id, defaultPhonemizerRegistry)
}
