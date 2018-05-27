package hangulize

import (
	"fmt"
	"strings"
)

// Subword is a chunk of a word with a level number.  The level indicates which
// pipeline step generated this subword.
type Subword struct {
	word  string
	level int
}

// -----------------------------------------------------------------------------

// SubwordsBuilder is a buffer to build a []Subword array.
type SubwordsBuilder struct {
	subwords []Subword
}

// NewSubwordsBuilder creates a SubwordsBuilder.
func NewSubwordsBuilder(subwords []Subword) *SubwordsBuilder {
	return &SubwordsBuilder{subwords}
}

// String() concatenates buffered subwords to assemble the full word.
func (b *SubwordsBuilder) String() string {
	var buf strings.Builder
	for _, subword := range b.subwords {
		buf.WriteString(subword.word)
	}
	return buf.String()
}

// Append extends the underlying subwords by the given ones.
func (b *SubwordsBuilder) Append(subwords ...Subword) {
	b.subwords = append(b.subwords, subwords...)
}

// Reset discards the underlying subwords.
func (b *SubwordsBuilder) Reset() {
	b.subwords = b.subwords[:0]
}

// Subwords builds the buffered subwords into a []Subword array.  It merges
// adjoin subwords if they share the same level.
func (b *SubwordsBuilder) Subwords() []Subword {
	subwords := make([]Subword, 0)

	if len(b.subwords) == 0 {
		// No subwords buffered.
		return subwords
	}

	// Merge same level adjoin subwords.
	var buf strings.Builder
	mergingLevel := -1

	for _, subword := range b.subwords {
		if subword.level != mergingLevel {
			// Keep the merged subword.
			merged := &Subword{buf.String(), mergingLevel}
			subwords = append(subwords, *merged)

			// Open a new one.
			buf.Reset()
		}

		buf.WriteString(subword.word)
		mergingLevel = subword.level
	}

	merged := &Subword{buf.String(), mergingLevel}
	subwords = append(subwords, *merged)

	return subwords
}

// -----------------------------------------------------------------------------

// Replacement is a deferred subword replacement.
type Replacement struct {
	start int
	stop  int
	word  string
}

func (r Replacement) String() string {
	return fmt.Sprintf(`[%d-%d] %#v`, r.start, r.stop, r.word)
}

// -----------------------------------------------------------------------------

// SubwordReplacer remembers replacements in a buffer.  Finally, it applies the
// replacements and splits the result in several subwords.
type SubwordReplacer struct {
	word string

	// Buffered replacements.
	repls []Replacement

	levels []int

	// The level for the replaced subwords.
	nextLevel int
}

// NewSubwordReplacer creates a SubwordReplacer for a word.
func NewSubwordReplacer(word string, prevLevel, nextLevel int) *SubwordReplacer {
	repls := make([]Replacement, 0)

	levels := make([]int, len(word))
	for i := 0; i < len(levels); i++ {
		levels[i] = prevLevel
	}

	return &SubwordReplacer{word, repls, levels, nextLevel}
}

// Replace buffers a replacement.
func (r *SubwordReplacer) Replace(start, stop int, word string) {
	r.ReplaceBy(Replacement{start, stop, word})
}

// ReplaceBy buffers multiple replacements.
func (r *SubwordReplacer) ReplaceBy(repls ...Replacement) {
	r.repls = append(r.repls, repls...)
}

// flush applies the buffered replacements to the SubwordReplacer internal.
func (r *SubwordReplacer) flush() {
	var buf strings.Builder
	levels := make([]int, 0)

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
func (r *SubwordReplacer) String() string {
	r.flush()
	return r.word
}

// Subwords applies the buffered replacements and returns the replaced word as
// a []Subword array.
func (r *SubwordReplacer) Subwords() []Subword {
	r.flush()

	subwords := make([]Subword, 0)

	if len(r.levels) == 0 {
		return subwords
	}

	level := r.levels[0]

	var buf strings.Builder

	for i, ch := range r.word {
		if r.levels[i] != level {
			subwords = append(subwords, Subword{buf.String(), level})
			level = r.levels[i]
			buf.Reset()
		}
		buf.WriteRune(ch)
	}
	subwords = append(subwords, Subword{buf.String(), level})

	return subwords
}
