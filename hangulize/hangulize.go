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

	// mask := h.newMask(word)
	// fmt.Println(word, mask)

	word = h.normalize(word, ch)

	chunks := []Chunk{Chunk{word, 0}}

	chunks = h.rewrite(chunks)
	fmt.Println("rewrite", chunks)

	chunks = h.transcribe(chunks)
	fmt.Println("transcribe", chunks)

	word = h.assembleJamo(chunks)
	fmt.Println("jamo", word)

	// word = NewChunkBuilder(chunks).String()
	// fmt.Println("join", word)

	// word = h.cleanUp(word)
	return word
}

// -----------------------------------------------------------------------------

// Rewrite applies multiple replacers on a word.
func (h *Hangulizer) rewrite(chunks []Chunk) []Chunk {
	var buf ChunkBuilder

	for _, chunk := range chunks {
		word := chunk.word
		age := chunk.age

		rep := NewChunkedReplacer(word, age, 0)

		for _, rule := range h.spec.rewrite {
			for _, r := range rule.Replacements(word) {
				rep.Replace(r.start, r.stop, r.word)
			}
			word = rep.String()
		}

		buf.Put(rep.Chunks())
	}

	return buf.Chunks()
}

func (h *Hangulizer) transcribe(chunks []Chunk) []Chunk {
	var buf ChunkBuilder

	for _, chunk := range chunks {
		word := chunk.word
		age := chunk.age

		rep := NewChunkedReplacer(word, age)
		dummy := NewChunkedReplacer(word, age)

		for _, rule := range h.spec.transcribe {
			for _, r := range rule.Replacements(word) {
				rep.Replace(r.start, r.stop, r.word)

				nulls := strings.Repeat("\x00", len(r.word))
				dummy.Replace(r.start, r.stop, nulls)
			}
			rep.flush()
			word = dummy.String()
		}

		buf.Put(rep.Chunks())
	}

	return buf.Chunks()
}

func (h *Hangulizer) assembleJamo(chunks []Chunk) string {
	var buf strings.Builder
	var jamoBuf strings.Builder

	for _, chunk := range chunks {
		// Don't touch age=0 chunks.
		// They have never rewritten or transcribed.
		if chunk.age == 0 {
			buf.WriteString(AssembleJamo(jamoBuf.String()))
			jamoBuf.Reset()

			buf.WriteString(chunk.word)
			continue
		}
		jamoBuf.WriteString(chunk.word)
	}
	buf.WriteString(jamoBuf.String())

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

func (h *Hangulizer) cleanUp(word string) string {
	markers := h.spec.Config.Markers

	args := make([]string, len(markers)*2)
	for i, marker := range markers {
		args[i*2] = string(marker)
		args[i*2+1] = ""
	}
	rep := strings.NewReplacer(args...)

	var buf strings.Builder

	for _, ch := range word {
		rep.WriteString(&buf, string(ch))
	}

	return buf.String()
}
