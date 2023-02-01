package hangulize

// Translit is an interface for a transliterator. It may convert a word from
// one script to another script. It also may guess phonograms from the spelling
// based on lexical analysis.
type Translit interface {
	// Scheme returns the identifier string of a Translit.
	Scheme() string

	// Transliterate transliterates the given word.
	Transliterate(string) (string, error)
}

// translitRegistry is a registry holding Translits.
type translitRegistry map[string]Translit

// Detach copies the registry as map[string]Translit.
func (r translitRegistry) Detach() map[string]Translit {
	copied := make(map[string]Translit, len(r))
	for m, t := range r {
		copied[m] = t
	}
	return copied
}

// Add registers a Translit into the registry.
func (r translitRegistry) Add(t Translit) bool {
	scheme := t.Scheme()

	if _, ok := r[scheme]; ok {
		// already exists
		return false
	}

	r[scheme] = t
	return true
}

// Remove deregisters a Translit from the registry.
func (r translitRegistry) Remove(scheme string) bool {
	_, ok := r[scheme]
	if ok {
		delete(r, scheme)
	}
	return ok
}

// defaultTranslitRegistry the default Translit registry.
var defaultTranslitRegistry = make(translitRegistry)

// Translits returns a copy of the default Translit registry.
func Translits() map[string]Translit {
	return defaultTranslitRegistry.Detach()
}

// UseTranslit imports a Translit into the default registry.
func UseTranslit(t Translit) bool {
	return defaultTranslitRegistry.Add(t)
}

// UnuseTranslit removes an imported Translit from the default registry.
func UnuseTranslit(scheme string) bool {
	return defaultTranslitRegistry.Remove(scheme)
}
