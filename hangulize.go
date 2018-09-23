package hangulize

// Hangulize transcribes a non-Korean word into Hangul, which is the Korean
// alphabet.
//
// For example, it will transcribe "Владивосто́к" in Russian into
// "블라디보스토크".
//
// It is the most simple and useful API of thie package.
//
func Hangulize(lang string, word string) string {
	spec, ok := LoadSpec(lang)
	if !ok {
		// spec not found
		return word
	}

	h := NewHangulizer(spec)
	return h.Hangulize(word)
}

// Hangulizer provides the transcription logic for the underlying spec.
type Hangulizer struct {
	spec        *Spec
	phonemizers map[string]Phonemizer
}

// NewHangulizer creates a Hangulizer for a spec.
func NewHangulizer(spec *Spec) *Hangulizer {
	return &Hangulizer{spec, make(map[string]Phonemizer)}
}

// Spec returns the underlying spec.
func (h *Hangulizer) Spec() *Spec {
	return h.spec
}

// UsePhonemizer keeps a phonemizer for ready to use.
func (h *Hangulizer) UsePhonemizer(p Phonemizer) bool {
	return usePhonemizer(p, h.phonemizers)
}

// UnusePhonemizer discards a phonemizer.
func (h *Hangulizer) UnusePhonemizer(id string) bool {
	return unusePhonemizer(id, h.phonemizers)
}

// GetPhonemizer returns a phonemizer by the ID.
func (h *Hangulizer) GetPhonemizer(id string) (Phonemizer, bool) {
	p, ok := getPhonemizer(id, h.phonemizers)
	return p, ok
}

// Hangulize transcribes a loanword into Hangul.
func (h *Hangulizer) Hangulize(word string) string {
	p := pipeline{h, nil}
	return p.forward(word)
}

// HangulizeTrace transcribes a loanword into Hangul
// and returns the traced internal events too.
func (h *Hangulizer) HangulizeTrace(word string) (string, Traces) {
	var tr tracer
	p := pipeline{h, &tr}

	word = p.forward(word)

	return word, tr.Traces()
}
