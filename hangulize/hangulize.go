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
	word = h.transcribe(word, ch)

	word = AssembleJamo(word, ch)
	word = h.removeMarkers(word, ch)

	return word
}

// -----------------------------------------------------------------------------

func (h *Hangulizer) normalize(word string, ch chan<- Trace) string {
	// TODO(sublee): Language-specific normalizer
	orig := word

	except := make([]string, 0)

	args := make([]string, 0)
	for to, froms := range h.spec.normalize {
		for _, from := range froms {
			args = append(args, from, to)
		}

		except = append(except, to)
	}
	rep := strings.NewReplacer(args...)
	word = rep.Replace(word)

	word = NormalizeRoman(word, except)

	word = strings.ToLower(word)

	trace(ch, word, orig, "roman")
	return word
}

func (h *Hangulizer) rewrite(word string, ch chan<- Trace) string {
	return h.spec.rewrite.Rewrite(word, ch)
}

func (h *Hangulizer) transcribe(word string, ch chan<- Trace) string {
	return h.spec.transcribe.Rewrite(word, ch)
}

func (h *Hangulizer) removeMarkers(word string, ch chan<- Trace) string {
	orig := word

	markers := h.spec.Config.Markers
	for _, marker := range markers {
		word = strings.Replace(word, string(marker), "", -1)
	}

	trace(ch, word, orig, "remove-markers")
	return word
}
