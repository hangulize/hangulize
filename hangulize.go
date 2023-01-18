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

	h := New(spec)
	return h.Hangulize(word)
}

// Hangulizer provides the transcription logic for the underlying spec.
type Hangulizer struct {
	spec               *Spec
	phonemizerRegistry map[string]Phonemizer
}

// New creates a Hangulizer for a spec.
func New(spec *Spec) *Hangulizer {
	return &Hangulizer{spec, make(map[string]Phonemizer)}
}

// NewHangulizer has been deprecated. Use New instead.
var NewHangulizer = New

// Spec returns the underlying spec.
func (h *Hangulizer) Spec() *Spec {
	return h.spec
}

// ImportPhonemizer keeps a phonemizer for ready to use.
func (h *Hangulizer) ImportPhonemizer(p Phonemizer) bool {
	return importPhonemizer(p, h.phonemizerRegistry)
}

// UnimportPhonemizer discards a phonemizer.
func (h *Hangulizer) UnimportPhonemizer(id string) bool {
	return unimportPhonemizer(id, h.phonemizerRegistry)
}

// Phonemizer returns a phonemizer by the ID.
func (h *Hangulizer) Phonemizer(id string) (Phonemizer, bool) {
	p, ok := h.phonemizerRegistry[id]
	return p, ok
}

// Hangulize transcribes a loanword into Hangul.
func (h *Hangulizer) Hangulize(word string) (string, error) {
	p := pipeline{h, nil}
	return p.forward(word)
}

// HangulizeTrace transcribes a loanword into Hangul
// and returns the traced internal events too.
func (h *Hangulizer) HangulizeTrace(word string) (string, Traces, error) {
	var tr tracer
	p := pipeline{h, &tr}

	word, err := p.forward(word)

	return word, tr.Traces(), err
}
