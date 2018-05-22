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

// Hangulizer ...
type Hangulizer struct {
	spec *Spec
}

// NewHangulizer ...
func NewHangulizer(spec *Spec) *Hangulizer {
	return &Hangulizer{spec}
}

// Hangulize ...
func (h *Hangulizer) Hangulize(word string) string {
	word = h.spec.Normalize(word)
	word = h.spec.Rewrite(word)
	word = h.spec.Hangulize(word)
	word = CompleteHangul(word)
	return word
}

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
