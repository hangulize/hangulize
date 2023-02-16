package subword

import (
	"bytes"
	"fmt"
)

// Replacement is a deferred sw Replacement.
type Replacement struct {
	Start int
	Stop  int
	Word  string
}

func (r Replacement) String() string {
	return fmt.Sprintf(`[%d-%d] %#v`, r.Start, r.Stop, r.Word)
}

// NewReplacement creates a Replacement.
func NewReplacement(start, stop int, word string) Replacement {
	return Replacement{start, stop, word}
}

// Replacer remembers replacements in a buffer. Finally, it applies the
// replacements and splits the result in several subwords.
type Replacer struct {
	word string

	// Buffered replacements.
	repls []Replacement

	levels []int

	// The level for the replaced subwords.
	nextLevel int
}

// NewReplacer creates a SubwordReplacer for a word.
func NewReplacer(word string, prevLevel, nextLevel int) *Replacer {
	var repls []Replacement

	levels := make([]int, len(word))
	for i := 0; i < len(levels); i++ {
		levels[i] = prevLevel
	}

	return &Replacer{word, repls, levels, nextLevel}
}

// Replace buffers a replacement.
func (r *Replacer) Replace(start, stop int, word string) {
	r.ReplaceBy(NewReplacement(start, stop, word))
}

// ReplaceBy buffers multiple replacements.
func (r *Replacer) ReplaceBy(repls ...Replacement) {
	r.repls = append(r.repls, repls...)
}

// commit applies the buffered replacements to the SubwordReplacer internal.
func (r *Replacer) commit() {
	var buf bytes.Buffer
	var levels []int

	offset := 0
	for _, repl := range r.repls {
		start := repl.Start
		stop := repl.Stop
		word := repl.Word

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
	r.repls = make([]Replacement, 0)
}

// String applies the buffered replacements and returns the replaced full word.
func (r *Replacer) String() string {
	r.commit()
	return r.word
}

// Subwords applies the buffered replacements and returns the replaced word as
// a []Subword array.
func (r *Replacer) Subwords() []Subword {
	r.commit()

	var subwords []Subword

	if len(r.levels) == 0 {
		return subwords
	}

	level := r.levels[0]

	var buf bytes.Buffer

	for i, ch := range r.word {
		if r.levels[i] != level {
			subwords = append(subwords, New(buf.String(), level))
			level = r.levels[i]
			buf.Reset()
		}
		buf.WriteRune(ch)
	}
	subwords = append(subwords, New(buf.String(), level))

	return subwords
}
