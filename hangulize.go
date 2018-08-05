package hangulize

// Hangulize is the most simple and useful API of thie package. It transcribes
// a non-Korean word into Hangul, which is the Korean alphabet. For example, it
// will transcribe "Владивосто́к" in Russian into "블라디보스토크".
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
	pronouncers map[string]Pronouncer
}

// NewHangulizer creates a Hangulizer for a spec.
func NewHangulizer(spec *Spec) *Hangulizer {
	return &Hangulizer{spec, make(map[string]Pronouncer)}
}

// Spec returns the underlying spec.
func (h *Hangulizer) Spec() *Spec {
	return h.spec
}

// UsePronouncer keeps a pronouncer for ready to use.
func (h *Hangulizer) UsePronouncer(p Pronouncer) bool {
	return usePronouncer(p, &h.pronouncers)
}

// UnusePronouncer discards a pronouncer.
func (h *Hangulizer) UnusePronouncer(id string) bool {
	return unusePronouncer(id, &h.pronouncers)
}

// GetPronouncer returns a pronouncer by the ID.
func (h *Hangulizer) GetPronouncer(id string) (Pronouncer, bool) {
	p, ok := getPronouncer(id, &h.pronouncers)
	if !ok {
		// Fallback by the global pronouncer registry.
		p, ok = getPronouncer(id, &globalPronouncers)
	}
	return p, ok
}

// Hangulize transcribes a loanword into Hangul.
func (h *Hangulizer) Hangulize(word string) string {
	p := pipeline{h, nil}
	return p.forward(word)
}

// HangulizeTrace transcribes a loanword into Hangul
// and returns the traced internal events too.
func (h *Hangulizer) HangulizeTrace(word string) (string, []Trace) {
	var tr tracer
	p := pipeline{h, &tr}

	word = p.forward(word)

	return word, tr.Traces()
}
