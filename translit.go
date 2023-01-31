package hangulize

// Translit is an interface for a transliterator. It may convert a word from
// one script to another script. It also may guess phonograms from the spelling
// based on lexical analysis.
type Translit interface {
	// Method returns the identifier string of a Translit.
	Method() string

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
	method := t.Method()

	if _, ok := r[method]; ok {
		// already exists
		return false
	}

	r[method] = t
	return true
}

// Remove deregisters a Translit from the registry.
func (r translitRegistry) Remove(method string) bool {
	_, ok := r[method]
	if ok {
		delete(r, method)
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
func UnuseTranslit(method string) bool {
	return defaultTranslitRegistry.Remove(method)
}
