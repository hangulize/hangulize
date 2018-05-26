package hangulize

import (
	"strings"
)

type Chunk struct {
	word  string
	level int
}

type ChunkBuilder struct {
	chunks []Chunk
}

func NewChunkBuilder(chunks []Chunk) *ChunkBuilder {
	return &ChunkBuilder{chunks}
}

func (b *ChunkBuilder) String() string {
	var buf strings.Builder
	for _, chunk := range b.chunks {
		buf.WriteString(chunk.word)
	}
	return buf.String()
}

func (b *ChunkBuilder) Put(chunks ...Chunk) {
	b.chunks = append(b.chunks, chunks...)
}

func (b *ChunkBuilder) Reset() {
	b.chunks = b.chunks[:0]
}

func (b *ChunkBuilder) Chunks() []Chunk {
	chunks := make([]Chunk, 0)

	if len(b.chunks) == 0 {
		return chunks
	}

	mergedChunk := Chunk{level: b.chunks[0].level}

	for _, chunk := range b.chunks {
		if chunk.level != mergedChunk.level {
			chunks = append(chunks, mergedChunk)
			mergedChunk = Chunk{level: chunk.level}
		}

		mergedChunk.word += chunk.word
	}

	chunks = append(chunks, mergedChunk)

	return chunks
}

// -----------------------------------------------------------------------------

// ChunkedReplacer remembers replacements in a buffer.  Finally, it applies the
// replacements and splits the result in several chunks.
type ChunkedReplacer struct {
	word      string
	nextLevel int

	levels []int
	repls  []Replacement
}

// NewChunkedReplacer creates a ChunkedReplacer for a word.
func NewChunkedReplacer(word string, level, nextLevel int) *ChunkedReplacer {
	levels := make([]int, len(word))
	for i := 0; i < len(levels); i++ {
		levels[i] = level
	}

	repls := make([]Replacement, 0)

	return &ChunkedReplacer{word, nextLevel, levels, repls}
}

// Replace remembers a replacement to be deferred.
func (r *ChunkedReplacer) Replace(start, stop int, word string) {
	r.repls = append(r.repls, Replacement{start, stop, word})
}

func (r *ChunkedReplacer) flush() {
	var buf strings.Builder
	levels := make([]int, 0)

	sortedRepls := make([]*Replacement, len(r.word))
	for i := range r.repls {
		repl := &r.repls[i]
		sortedRepls[repl.start] = repl
	}

	offset := 0
	for _, repl := range sortedRepls {
		if repl == nil {
			continue
		}

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

func (r *ChunkedReplacer) String() string {
	r.flush()
	return r.word
}

// Chunks returns [{subWord, isRewritten}...].
func (r *ChunkedReplacer) Chunks() []Chunk {
	r.flush()

	chunks := make([]Chunk, 0)

	if len(r.levels) == 0 {
		return chunks
	}

	level := r.levels[0]

	var buf strings.Builder

	for i, ch := range r.word {
		if r.levels[i] != level {
			chunks = append(chunks, Chunk{buf.String(), level})
			level = r.levels[i]
			buf.Reset()
		}
		buf.WriteRune(ch)
	}
	chunks = append(chunks, Chunk{buf.String(), level})

	return chunks
}
