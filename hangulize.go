package hangulize

// Hangulize transcribes a non-Korean word into Hangul, which is the Korean
// alphabet.
//
// For example, it transcribes "Владивосто́к" in Russian into "블라디보스토크".
//
// Also, this function is the most simple and useful API in this package.
func Hangulize(lang string, word string) (string, error) {
	spec, err := LoadSpec(lang)
	if err != nil {
		return word, err
	}

	h := &hangulizer{spec, defaultTranslitRegistry, nil}
	return h.Hangulize(word)
}

// Hangulizer is a transcriptor into Hangul dedicated for a specific language.
type Hangulizer interface {
	// Spec returns the underlying Spec.
	Spec() *Spec

	// Translits returns the imported Translits.
	Translits() map[string]Translit

	// UseTranslit imports a Translit.
	UseTranslit(Translit) bool

	// UnuseTranslit removes an imported a Translit.
	UnuseTranslit(scheme string) bool

	// Trace registers a tracing function.
	Trace(func(Trace))

	// Hangulize transcribes a non-Korean word into Hangul.
	Hangulize(word string) (string, error)
}

// hangulizer provides the transcription logic for the underlying spec.
type hangulizer struct {
	spec             *Spec
	translitRegistry translitRegistry
	traceFunc        func(Trace)
}

// New creates a hangulizer for a Spec.
func New(spec *Spec) Hangulizer {
	return &hangulizer{spec, make(translitRegistry), nil}
}

// Spec returns the underlying Spec.
func (h *hangulizer) Spec() *Spec {
	return h.spec
}

// Translits returns the imported Translits.
func (h *hangulizer) Translits() map[string]Translit {
	return h.translitRegistry.Detach()
}

// UseTranslit imports a Translit.
func (h *hangulizer) UseTranslit(t Translit) bool {
	return h.translitRegistry.Add(t)
}

// UnuseTranslit removes an imported Translit.
func (h *hangulizer) UnuseTranslit(scheme string) bool {
	return h.translitRegistry.Remove(scheme)
}

// Trace registers a tracing function.
func (h *hangulizer) Trace(fn func(Trace)) {
	h.traceFunc = fn
}

// Hangulize transcribes a non-Korean word into Hangul.
func (h *hangulizer) Hangulize(word string) (string, error) {
	p := newProcedure(h.Spec(), h.Translits(), h.traceFunc)
	return p.forward(word)
}
