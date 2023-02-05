package hangulize

import "fmt"

// Hangulize transcribes a non-Korean word into Hangul, which is the Korean
// alphabet.
//
// For example, it will transcribe "Владивосто́к" in Russian into
// "블라디보스토크".
//
// It is the most simple and useful API of thie package.
func Hangulize(lang string, word string) (string, error) {
	spec, ok := LoadSpec(lang)
	if !ok {
		// spec not found
		return word, fmt.Errorf("%w: %s", ErrSpecNotFound, lang)
	}

	h := &hangulizer{spec, defaultTranslitRegistry, nil}
	return h.Hangulize(word)
}

// Hangulizer is dedicated for a specific language. transcribes a provides the transcription logic for the underlying spec.
type Hangulizer interface {
	Spec() *Spec

	Translits() map[string]Translit
	UseTranslit(Translit) bool
	UnuseTranslit(scheme string) bool

	Trace(func(Trace))

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

// Translits returns the registered Translits.
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

// Spec returns the underlying Spec.
func (h *hangulizer) Trace(fn func(Trace)) {
	h.traceFunc = fn
}

// Hangulize transcribes a loanword into Hangul.
func (h *hangulizer) Hangulize(word string) (string, error) {
	p := newProcedure(h.Spec(), h.Translits(), h.traceFunc)
	return p.forward(word)
}
