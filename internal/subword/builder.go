package subword

import "bytes"

// Builder is a buffer to build a []Subword array.
type Builder struct {
	subwords []Subword
}

// NewBuilder creates a Builder from subwords.
func NewBuilder(subwords []Subword) *Builder {
	return &Builder{subwords}
}

// String() concatenates buffered subwords to assemble the full word.
func (b *Builder) String() string {
	var buf bytes.Buffer
	for _, sw := range b.subwords {
		buf.WriteString(sw.Word)
	}
	return buf.String()
}

// Write extends the underlying subwords by the given ones.
func (b *Builder) Write(subwords ...Subword) {
	b.subwords = append(b.subwords, subwords...)
}

// Reset discards the underlying subwords.
func (b *Builder) Reset() {
	b.subwords = b.subwords[:0]
}

// Subwords builds the buffered subwords into a []Subword array. It merges
// adjoin subwords if they share the same level.
func (b *Builder) Subwords() []Subword {
	var subwords []Subword

	if len(b.subwords) == 0 {
		// No subwords buffered.
		return subwords
	}

	// Merge same level adjoin subwords.
	var buf bytes.Buffer
	mergingLevel := -1

	for _, sw := range b.subwords {
		if sw.Level != mergingLevel && mergingLevel != -1 {
			// Keep the merged sw.
			merged := New(buf.String(), mergingLevel)
			subwords = append(subwords, merged)

			// Open a new one.
			buf.Reset()
		}

		buf.WriteString(sw.Word)
		mergingLevel = sw.Level
	}

	merged := New(buf.String(), mergingLevel)
	subwords = append(subwords, merged)

	return subwords
}
