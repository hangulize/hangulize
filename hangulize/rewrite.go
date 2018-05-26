package hangulize

import (
	"strings"
)

type Chunk struct {
	word string
	age  int
}

// Rewriter remembers replacements in a buffer.  Finally, it applies the
// replacements and splits the result in several chunks.
type Rewriter struct {
	word  string
	ages  []int
	repls []Replacement
}

// NewRewriter creates a Rewriter for a word.
func NewRewriter(word string, age int) *Rewriter {
	ages := make([]int, len(word))
	if age != 0 {
		for i := 0; i < len(ages); i++ {
			ages[i] = age
		}
	}

	repls := make([]Replacement, 0)

	return &Rewriter{word, ages, repls}
}

// Rewrite remembers a replacement.
func (r *Rewriter) Rewrite(repl Replacement) {
	r.repls = append(r.repls, repl)
}

func (r *Rewriter) flush() {
	var buf strings.Builder
	ages := make([]int, 0)

	offset := 0
	for _, repl := range r.repls {
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

func (r *Rewriter) String() string {
	r.flush()
	return r.word
}

// Chunks returns [{subWord, isRewritten}...].
func (r *Rewriter) Chunks() []Chunk {
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
func Rewrite(word string, rules []*Rule, age int) []Chunk {
	rewr := NewRewriter(word, age)

	for _, rule := range rules {
		for _, r := range rule.Replacements(word) {
			rewr.Rewrite(r)
		}
		word = rewr.String()
	}

	return rewr.Chunks()
}

func mergeChunks(chunks []Chunk) []Chunk {
	buf := make([]Chunk, 0)

	mergedChunk := Chunk{age: chunks[0].age}

	for _, chunk := range chunks {
		if chunk.age != mergedChunk.age {
			buf = append(buf, mergedChunk)
			mergedChunk = Chunk{age: chunk.age}
		}

		mergedChunk.word += chunk.word
	}
	buf = append(buf, mergedChunk)

	return buf
}

func RewriteChunks(chunks []Chunk, rules []*Rule, age int) []Chunk {
	buf := make([]Chunk, 0)

	for _, chunk := range chunks {
		if chunk.age == age {
			buf = append(buf, Rewrite(chunk.word, rules, chunk.age)...)
		} else {
			buf = append(buf, chunk)
		}
	}

	return mergeChunks(buf)
}

func CleanUpChunks(chunks []Chunk, age int) []Chunk {
	buf := make([]Chunk, 0)

	for _, chunk := range chunks {
		if chunk.age != age {
			buf = append(buf, chunk)
		}
	}

	return mergeChunks(buf)
}

func JoinChunks(chunks []Chunk) string {
	var buf strings.Builder
	for _, chunk := range chunks {
		buf.WriteString(chunk.word)
	}
	return buf.String()
}

// const done = "\x00"

// func Transcribe(word string, rules []*Rule) string {
// 	transcribed := make([]string, len(word))

// 	for _, rule := range rules {
// 		var buf strings.Builder

// 		offset := 0
// 		for _, r := range rule.Replacements(word) {
// 			buf.WriteString(word[offset:r.start])
// 			buf.WriteString(done)
// 			transcribed[r.start] = r.word
// 			offset = r.stop
// 		}
// 		buf.WriteString(word[offset:])

// 		word = buf.String()
// 	}

// 	for i, ch := range word {
// 		fmt.Println(i, string(ch))
// 		if string(ch) != done {
// 			transcribed[i] = string(ch)
// 		}
// 	}

// 	return strings.Join(transcribed, "")
// }
