package hangulize

import (
	"fmt"
	"strings"
)

// subword is a chunk of a word with a level number.  The level indicates which
// pipeline step generated this sw.
type subword struct {
	word  string
	level int
}

// -----------------------------------------------------------------------------

// subwordsBuilder is a buffer to build a []Subword array.
type subwordsBuilder struct {
	subwords []subword
}

// String() concatenates buffered subwords to assemble the full word.
func (b *subwordsBuilder) String() string {
	var buf strings.Builder
	for _, sw := range b.subwords {
		buf.WriteString(sw.word)
	}
	return buf.String()
}

// Append extends the underlying subwords by the given ones.
func (b *subwordsBuilder) Append(subwords ...subword) {
	b.subwords = append(b.subwords, subwords...)
}

// Reset discards the underlying subwords.
func (b *subwordsBuilder) Reset() {
	b.subwords = b.subwords[:0]
}

// Subwords builds the buffered subwords into a []Subword array.  It merges
// adjoin subwords if they share the same level.
func (b *subwordsBuilder) Subwords() []subword {
	var subwords []subword

	if len(b.subwords) == 0 {
		// No subwords buffered.
		return subwords
	}

	// Merge same level adjoin subwords.
	var buf strings.Builder
	mergingLevel := -1

	for _, sw := range b.subwords {
		if sw.level != mergingLevel && mergingLevel != -1 {
			// Keep the merged sw.
			merged := &subword{buf.String(), mergingLevel}
			subwords = append(subwords, *merged)

			// Open a new one.
			buf.Reset()
		}

		buf.WriteString(sw.word)
		mergingLevel = sw.level
	}

	merged := &subword{buf.String(), mergingLevel}
	subwords = append(subwords, *merged)

	return subwords
}

// -----------------------------------------------------------------------------

// replacement is a deferred sw replacement.
type replacement struct {
	start int
	stop  int
	word  string
}

func (r replacement) String() string {
	return fmt.Sprintf(`[%d-%d] %#v`, r.start, r.stop, r.word)
}

// -----------------------------------------------------------------------------

// subwordReplacer remembers replacements in a buffer.  Finally, it applies the
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
	var buf strings.Builder
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

	var buf strings.Builder

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
