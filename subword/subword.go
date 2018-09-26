package subword

// Subword is a chunk of a word with a level number. The level indicates which
// pipeline step generated this sw.
type Subword struct {
	Word  string
	Level int
}

// New creates a Subword.
func New(word string, level int) Subword {
	return Subword{word, level}
}
