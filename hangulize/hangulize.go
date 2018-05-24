/*
Package hangulize provides an automatic transcriber into Hangul for non-Korean
words.

Originally, Hangulize was implemented in Python in 2010.  This implementation
is a reboot with attractive feature improvements.

Original implementation: https://github.com/sublee/hangulize

Brian Jongseong Park proposed the seed idea of Hangulize on his Blog.

Post by Brian: http://iceager.egloos.com/2610028

*/
package hangulize

import (
	"strings"
)

// Hangulize transcribes a non-Korean word into Hangul, the Korean alphabet:
//
//    Hangulize("ita", "gloria")
//    // Output: "글로리아"
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

// -----------------------------------------------------------------------------

// Hangulizer ...
type Hangulizer struct {
	spec *Spec
}

// NewHangulizer ...
func NewHangulizer(spec *Spec) *Hangulizer {
	return &Hangulizer{spec}
}

// Hangulize transcribes a loanword into Hangul.
func (h *Hangulizer) Hangulize(word string) string {
	return h.HangulizeTrace(word, nil)
}

// HangulizeTrace transcribes a loanword into Hangul.  During
// transcribing, it sends internal traces to the given channel.
func (h *Hangulizer) HangulizeTrace(word string, ch chan<- Trace) string {
	if ch != nil {
		defer close(ch)
	}
	trace(ch, word, "", "input")

	word = h.normalize(word, ch)
	word = h.rewrite(word, ch)
	word = h.hangulize(word, ch)
	// word = h.spec._RemoveMarkers(word, ch)

	word = _CompleteHangul(word, ch)
	return word
}

// -----------------------------------------------------------------------------

func (h *Hangulizer) normalize(word string, ch chan<- Trace) string {
	// TODO(sublee): Language-specific normalizer
	orig := word
	word = strings.ToLower(word)
	trace(ch, word, orig, "to-lower")
	return word
}

func (h *Hangulizer) rewrite(word string, ch chan<- Trace) string {
	return h.spec.rewrite.Rewrite(word, ch)
}

func (h *Hangulizer) hangulize(word string, ch chan<- Trace) string {
	return h.spec.hangulize.Rewrite(word, ch)
}
