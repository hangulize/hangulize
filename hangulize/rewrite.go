package hangulize

import (
	"strings"
)

type Chunk struct {
	word string
	age  int
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

func (b *ChunkBuilder) Put(chunks []Chunk) {
	b.chunks = append(b.chunks, chunks...)
}

func (b *ChunkBuilder) Chunks() []Chunk {
	chunks := make([]Chunk, 0)

	if len(b.chunks) == 0 {
		return chunks
	}

	mergedChunk := Chunk{age: b.chunks[0].age}

	for _, chunk := range b.chunks {
		if chunk.age != mergedChunk.age {
			chunks = append(chunks, mergedChunk)
			mergedChunk = Chunk{age: chunk.age}
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
	word  string
	ages  []int
	repls []Replacement
}

// NewChunkedReplacer creates a ChunkedReplacer for a word.
func NewChunkedReplacer(word string, age int) *ChunkedReplacer {
	ages := make([]int, len(word))
	if age != 0 {
		for i := 0; i < len(ages); i++ {
			ages[i] = age
		}
	}

	repls := make([]Replacement, 0)

	return &ChunkedReplacer{word, ages, repls}
}

// Replace remembers a replacement to be deferred.
func (r *ChunkedReplacer) Replace(start, stop int, word string) {
	r.repls = append(r.repls, Replacement{start, stop, word})
}

func (r *ChunkedReplacer) flush() {
	var buf strings.Builder
	ages := make([]int, 0)

	sortedRepls := make([]*Replacement, len(r.word))
	for i := range r.repls {
		repl := &r.repls[i]
		fmt.Println(sortedRepls, r.word, repl)
		sortedRepls[repl.start] = repl
	}

	offset := 0
	for _, repl := range sortedRepls {
		if repl == nil {
			continue
		}

		// before replacement
		buf.WriteString(r.word[offset:repl.start])
		ages = append(ages, r.ages[offset:repl.start]...)

		// replacement
		buf.WriteString(repl.word)
		for i := 0; i < len(repl.word); i++ {
			prevAge := r.ages[repl.start]
			ages = append(ages, prevAge+1)
		}

		offset = repl.stop
	}
	// after replacement
	buf.WriteString(r.word[offset:])
	ages = append(ages, r.ages[offset:]...)

	r.word = buf.String()
	r.ages = ages
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

	if len(r.ages) == 0 {
		return chunks
	}

	age := r.ages[0]

	var buf strings.Builder

	for i, ch := range r.word {
		if r.ages[i] != age {
			chunks = append(chunks, Chunk{buf.String(), age})
			age = r.ages[i]
			buf.Reset()
		}
		buf.WriteRune(ch)
	}
	chunks = append(chunks, Chunk{buf.String(), age})

	return chunks
}

// -----------------------------------------------------------------------------

// Rewrite applies multiple replacers on a word.
func Rewrite(chunks []Chunk, rules []*Rule) []Chunk {
	var buf ChunkBuilder

	for _, chunk := range chunks {
		word := chunk.word
		age := chunk.age

		rep := NewChunkedReplacer(word, age)

		for _, rule := range rules {
			for _, r := range rule.Replacements(word) {
				rep.Replace(r.start, r.stop, r.word)
			}
			word = rep.String()
		}

		buf.Put(rep.Chunks())
	}

	return buf.Chunks()
}

func Replace(chunks []Chunk, rules []*Rule) []Chunk {
	var buf ChunkBuilder

	for _, chunk := range chunks {
		word := chunk.word
		age := chunk.age

		rep := NewChunkedReplacer(word, age)
		dummy := NewChunkedReplacer(word, age)

		for _, rule := range rules {
			for _, r := range rule.Replacements(word) {
				fmt.Printf("%#v %#v %#v\n", word, r.start, r.stop)
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
