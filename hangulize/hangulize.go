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

	chunks := h.begin(word)
	word1 := NewChunkBuilder(chunks).String()
	trace(ch, word1, word, "start")

	chunks = h.rewrite(chunks)
	word2 := NewChunkBuilder(chunks).String()
	trace(ch, word2, word1, "rewrite")

	chunks = h.transcribe(chunks)
	word3 := NewChunkBuilder(chunks).String()
	trace(ch, word3, word2, "transcribe")

	word = h.assembleJamo(chunks)
	trace(ch, word, word3, "jamo")

	return word
}

// -----------------------------------------------------------------------------
// Hangulize pipeline

func (h *Hangulizer) begin(word string) []Chunk {
	// Detect used letters.
	var buf strings.Builder
	rules := append(h.spec.rewrite, h.spec.transcribe...)
	for _, rule := range rules {
		buf.WriteString(regexpLetters(rule.from.re.String()))
		buf.WriteString(regexpLetters(rule.from.neg.String()))
	}

	// Choose non-marker letters.
	markers := set(h.spec.Config.Markers)

	letters := make([]string, 0)
	for _, ch := range buf.String() {
		let := string(ch)
		if inSet(let, markers) {
			continue
		}
		letters = append(letters, let)
	}

	letters = set(letters)

	// Split the word by their letters.
	rep := NewChunkedReplacer(word, 0, 1)

	for i, ch := range word {
		let := string(ch)
		if inSet(let, letters) {
			rep.Replace(i, i+len(let), let)
		}
	}

	return rep.Chunks()
}

// Rewrite applies multiple replacers on a word.
func (h *Hangulizer) rewrite(chunks []Chunk) []Chunk {
	var buf ChunkBuilder

	for _, chunk := range chunks {
		word := chunk.word
		level := chunk.level

		rep := NewChunkedReplacer(word, level, 1)

		for _, rule := range h.spec.rewrite {
			for _, r := range rule.Replacements(word) {
				rep.Replace(r.start, r.stop, r.word)
			}
			word = rep.String()
		}

		buf.Put(rep.Chunks()...)
	}

	return buf.Chunks()
}

func (h *Hangulizer) transcribe(chunks []Chunk) []Chunk {
	var buf ChunkBuilder

	for _, chunk := range chunks {
		word := chunk.word
		level := chunk.level

		rep := NewChunkedReplacer(word, level, 2)
		dummy := NewChunkedReplacer(word, 0, 0)

		for _, rule := range h.spec.transcribe {
			for _, r := range rule.Replacements(word) {
				rep.Replace(r.start, r.stop, r.word)

				nulls := strings.Repeat("\x00", len(r.word))
				dummy.Replace(r.start, r.stop, nulls)
			}
			rep.flush()
			word = dummy.String()
		}

		buf.Put(rep.Chunks()...)
	}

	chunks = buf.Chunks()
	buf.Reset()

	for _, chunk := range chunks {
		if chunk.level == 1 {
			continue
		}
		buf.Put(chunk)
	}

	return buf.Chunks()
}

func (h *Hangulizer) assembleJamo(chunks []Chunk) string {
	var buf strings.Builder
	var jamoBuf strings.Builder

	for _, chunk := range chunks {
		// Don't touch age=0 chunks.
		// They have never rewritten or transcribed.
		if chunk.level == 0 {
			buf.WriteString(AssembleJamo(jamoBuf.String()))
			jamoBuf.Reset()

			buf.WriteString(chunk.word)
			continue
		}
		jamoBuf.WriteString(chunk.word)
	}
	buf.WriteString(AssembleJamo(jamoBuf.String()))

	return buf.String()
}

// -----------------------------------------------------------------------------

func (h *Hangulizer) newMask(word string) []bool {
	text := []rune(word)
	return make([]bool, len(text))
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

	word = Normalize(word, RomanNormalizer{}, except)

	word = strings.ToLower(word)

	trace(ch, word, orig, "roman")
	return word
}
