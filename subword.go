package hangulize

import (
	"bytes"
)

// subword is a chunk of a word with a level number. The level indicates which
// pipeline step generated this sw.
type subword struct {
	word  string
	level int
}

// subwordsBuilder is a buffer to build a []Subword array.
type subwordsBuilder struct {
	subwords []subword
}

// String() concatenates buffered subwords to assemble the full word.
func (b *subwordsBuilder) String() string {
	var buf bytes.Buffer
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

// Subwords builds the buffered subwords into a []Subword array. It merges
// adjoin subwords if they share the same level.
func (b *subwordsBuilder) Subwords() []subword {
	var subwords []subword

	if len(b.subwords) == 0 {
		// No subwords buffered.
		return subwords
	}

	// Merge same level adjoin subwords.
	var buf bytes.Buffer
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
