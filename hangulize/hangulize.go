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

	word = h.composeHangul(subwords)
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
// The Hangulize Pipeline

// 1. Normalize (Word -> Word)
//
// This step eliminates letter case to make the next steps work easier.
//
// For example, "Hello" in Roman script will be normalized to "hello".
//
func (h *Hangulizer) normalize(word string) string {
	word = h.spec.normReplacer.Replace(word)

	norm := h.spec.norm
	except := h.spec.normLetters
	word = Normalize(word, norm, except)

	return word
}

// 2. Group meaningful letters (Word -> Subwords[level=0 or 1])
//
// Meaningful letter is the letter which appears in the rewrite/transcribe
// rules.  This step groups letters by their meaningness into subwords.  A
// meaningful subword has level=1 meanwhile meaningless one has level=0.
//
// For example, "hello, world!" will be grouped into
// [{"hello",1}, {", ",0}, {"world",1}, {"!",0}].
//
func (h *Hangulizer) draft(word string) []Subword {
	rep := NewSubwordReplacer(word, 0, 1)

	for i, ch := range word {
		let := string(ch)
		if inSet(let, h.spec.letters) {
			rep.Replace(i, i+len(let), let)
		}
	}

	return rep.Subwords()
}

// 3. Rewrite (Subwords -> Subwords[level=1])
//
// This step minimizes the gap between pronunciation and spelling.
//
// It replaces the word by a rule.  The replaced word will be used as the input
// for the next rule.  Each result becomes the next input.  That's why this
// step called "rewrite".
//
// For example, "hello" can be rewritten to "heˈlō".
//
// NOTE(sublee): But this step has a limitation.  It guesses a pronunciation
// from the spelling.  But it can be too hard for some script systems, such as
// English or Franch.
//
func (h *Hangulizer) rewrite(subwords []Subword) []Subword {
	var swBuf SubwordsBuilder

	for _, subword := range subwords {
		word := subword.word
		level := subword.level

		rep := NewSubwordReplacer(word, level, 1)

		for _, rule := range h.spec.Rewrite {
			repls := rule.Replacements(word)
			rep.ReplaceBy(repls...)
			word = rep.String()
		}

		swBuf.Append(rep.Subwords()...)
	}

	return swBuf.Subwords()
}

// 4. Transcribe (Subwords -> Subwords[level=2])
//
// This step determines Hangul spelling for the pronunciation.
//
// Rather than generating composed Hangul syllables, it enumerates decomposed
// Jamo phonemes, such as "ㅎㅏ-ㄴ".  In this form, a Jaeum after a hyphen
// ("-ㄴ") means that it is a Jongseong (tail).
//
// For example, "heˈlō" can be transcribed as "ㅎㅔ-ㄹㄹㅗ".
//
func (h *Hangulizer) transcribe(subwords []Subword) []Subword {
	var swBuf SubwordsBuilder

	for _, subword := range subwords {
		word := subword.word
		level := subword.level

		rep := NewSubwordReplacer(word, level, 2)

		// transcribe is not rewrite.  A result of a replacement is not the
		// input of the next replacement.  dummy marks the replaced subwords
		// with NULL characters.
		dummy := NewSubwordReplacer(word, 0, 0)

		for _, rule := range h.spec.Transcribe {
			repls := rule.Replacements(word)
			rep.ReplaceBy(repls...)

			for _, repl := range repls {
				nulls := strings.Repeat("\x00", len(repl.word))
				dummy.Replace(repl.start, repl.stop, nulls)
			}

			rep.flush()
			word = dummy.String()
		}

		swBuf.Append(rep.Subwords()...)
	}

	// Discard level=1 subwords.  They have been generated by "3. Rewrite" but
	// never transcribed.  They are superfluity of the internal behavior.
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

// 5. Compose Hangul (Subwords -> Word)
//
// This step converts decomposed Jamo phonemes to composed Hangul syllables.
//
// For example, "ㅎㅔ-ㄹㄹㅗ" will be "헬로".
//
func (h *Hangulizer) composeHangul(subwords []Subword) string {
	var buf strings.Builder
	var jamoBuf strings.Builder

	for _, subword := range subwords {
		// Don't touch level=0 subwords.  They just have passed through the
		// pipeline, because they are meaningless.
		if subword.level == 0 {
			buf.WriteString(ComposeHangul(jamoBuf.String()))
			jamoBuf.Reset()

			buf.WriteString(subword.word)
			continue
		}
		jamoBuf.WriteString(subword.word)
	}
	buf.WriteString(ComposeHangul(jamoBuf.String()))

	return buf.String()
}
