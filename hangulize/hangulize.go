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
	"fmt"
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
	// ch   chan<- Trace
}

// NewHangulizer ...
func NewHangulizer(spec *Spec /*, ch chan<- Trace*/) *Hangulizer {
	return &Hangulizer{spec /*, ch*/}
}

// Hangulize transcribes a loanword into Hangul.
func (h *Hangulizer) Hangulize(word string) string {
	// trace(h.ch, word, "", "input")

	word = h.normalize(word)

	subwords := h.draft(word)
	// word1 := NewSubwordsBuilder(subwords).String()
	// trace(ch, word1, word, "start")

	subwords = h.rewrite(subwords)
	// word2 := NewSubwordsBuilder(subwords).String()
	// trace(ch, word2, word1, "rewrite")

	subwords = h.transcribe(subwords)
	// word3 := NewSubwordsBuilder(subwords).String()
	// trace(ch, word3, word2, "transcribe")

	word = h.assembleJamo(subwords)
	// trace(ch, word, word3, "jamo")

	return word
}

// Trace is emitted when a replacement occurs.  It is used for tracing of
// Hangulize pipeline internal.
type Trace struct {
	why  string
	from string
	to   string
}

func (e *Trace) String() string {
	return fmt.Sprintf("%#v %s", e.to, e.why)
}

func (h *Hangulizer) trace() {
	// if h.ch == nil {
	// 	return
	// }
}

// -----------------------------------------------------------------------------
// Hangulize pipeline

// 1. Normalize input word.
func (h *Hangulizer) normalize(word string) string {
	// TODO(sublee): Language-specific normalizer
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

	word = Normalize(word, RomanNormalizer{}, except)

	word = strings.ToLower(word)

	return word
}

func (h *Hangulizer) draft(word string) []Subword {
	// Detect letters used in patterns except markers.
	rules := append(h.spec.rewrite, h.spec.transcribe...)
	markers := set(h.spec.Config.Markers)

	letters := make([]string, 0)

	for _, rule := range rules {
		for _, let := range rule.from.letters {
			if inSet(let, markers) {
				continue
			}
			letters = append(letters, let)
		}
	}

	letters = set(letters)

	// Split the word by their letters.
	rep := NewSubwordReplacer(word, 0, 1)

	for i, ch := range word {
		let := string(ch)
		if inSet(let, letters) {
			rep.Replace(i, i+len(let), let)
		}
	}

	return rep.Subwords()
}

// Rewrite applies multiple replacers on a word.
func (h *Hangulizer) rewrite(subwords []Subword) []Subword {
	var swBuf SubwordsBuilder

	for _, subword := range subwords {
		word := subword.word
		level := subword.level

		rep := NewSubwordReplacer(word, level, 1)

		for _, rule := range h.spec.rewrite {
			for _, r := range rule.Replacements(word) {
				rep.Replace(r.start, r.stop, r.word)
			}
			word = rep.String()
		}

		swBuf.Append(rep.Subwords()...)
	}

	return swBuf.Subwords()
}

func (h *Hangulizer) transcribe(subwords []Subword) []Subword {
	var swBuf SubwordsBuilder

	for _, subword := range subwords {
		word := subword.word
		level := subword.level

		rep := NewSubwordReplacer(word, level, 2)
		dummy := NewSubwordReplacer(word, 0, 0)

		for _, rule := range h.spec.transcribe {
			for _, r := range rule.Replacements(word) {
				rep.Replace(r.start, r.stop, r.word)

				nulls := strings.Repeat("\x00", len(r.word))
				dummy.Replace(r.start, r.stop, nulls)
			}
			rep.flush()
			word = dummy.String()
		}

		swBuf.Append(rep.Subwords()...)
	}

	subwords = swBuf.Subwords()
	swBuf.Reset()

	for _, subword := range subwords {
		if subword.level == 1 {
			continue
		}
		swBuf.Append(subword)
	}

	return swBuf.Subwords()
}

func (h *Hangulizer) assembleJamo(subwords []Subword) string {
	var buf strings.Builder
	var jamoBuf strings.Builder

	for _, subword := range subwords {
		// Don't touch age=0 subwords.
		// They have never rewritten or transcribed.
		if subword.level == 0 {
			buf.WriteString(AssembleJamo(jamoBuf.String()))
			jamoBuf.Reset()

			buf.WriteString(subword.word)
			continue
		}
		jamoBuf.WriteString(subword.word)
	}
	buf.WriteString(AssembleJamo(jamoBuf.String()))

	return buf.String()
}
