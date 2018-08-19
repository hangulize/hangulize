package hangulize

import (
	"bytes"
	"fmt"
)

// replacement is a deferred sw replacement.
type replacement struct {
	start int
	stop  int
	word  string
}

func (r replacement) String() string {
	return fmt.Sprintf(`[%d-%d] %#v`, r.start, r.stop, r.word)
}

// subwordReplacer remembers replacements in a buffer. Finally, it applies the
// replacements and splits the result in several subwords.
type subwordReplacer struct {
	word string

	// Buffered replacements.
	repls []replacement

	levels []int

	// The level for the replaced subwords.
	nextLevel int
}

// newSubwordReplacer creates a SubwordReplacer for a word.
func newSubwordReplacer(word string, prevLevel, nextLevel int) *subwordReplacer {
	var repls []replacement

	levels := make([]int, len(word))
	for i := 0; i < len(levels); i++ {
		levels[i] = prevLevel
	}

	return &subwordReplacer{word, repls, levels, nextLevel}
}

// Replace buffers a replacement.
func (r *subwordReplacer) Replace(start, stop int, word string) {
	r.ReplaceBy(replacement{start, stop, word})
}

// ReplaceBy buffers multiple replacements.
func (r *subwordReplacer) ReplaceBy(repls ...replacement) {
	r.repls = append(r.repls, repls...)
}

// flush applies the buffered replacements to the SubwordReplacer internal.
func (r *subwordReplacer) flush() {
	var buf bytes.Buffer
	var levels []int

	offset := 0
	for _, repl := range r.repls {
		start := repl.start
		stop := repl.stop
		word := repl.word

		// before replacement
		buf.WriteString(r.word[offset:start])
		levels = append(levels, r.levels[offset:start]...)

		// replacement
		buf.WriteString(word)
		for i := 0; i < len(word); i++ {
			levels = append(levels, r.nextLevel)
		}

		offset = stop
	}
	// after replacement
	buf.WriteString(r.word[offset:])
	levels = append(levels, r.levels[offset:]...)

	r.word = buf.String()
	r.levels = levels
	r.repls = r.repls[:0]
}

// String applies the buffered replacements and returns the replaced full word.
func (r *subwordReplacer) String() string {
	r.flush()
	return r.word
}

// Subwords applies the buffered replacements and returns the replaced word as
// a []Subword array.
func (r *subwordReplacer) Subwords() []subword {
	r.flush()

	var subwords []subword

	if len(r.levels) == 0 {
		return subwords
	}

	level := r.levels[0]

	var buf bytes.Buffer

	for i, ch := range r.word {
		if r.levels[i] != level {
			subwords = append(subwords, subword{buf.String(), level})
			level = r.levels[i]
			buf.Reset()
		}
		buf.WriteRune(ch)
	}
	subwords = append(subwords, subword{buf.String(), level})

	return subwords
}
