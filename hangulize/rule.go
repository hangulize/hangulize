package hangulize

import "fmt"

// Replacement ...
type Replacement struct {
	start int
	stop  int
	word  string
}

func (r Replacement) String() string {
	return fmt.Sprintf(`%d-%d -> %s`, r.start, r.stop, r.word)
}

// Rule is a replacer based on Pattern and RPatterns.
type Rule struct {
	from *Pattern
	to   *RPattern
}

// Replacements indicates which ranges should be replaced.
func (r *Rule) Replacements(word string) []Replacement {
	rs := make([]Replacement, 0)

	for _, m := range r.from.Find(word, -1) {
		start, stop := m[0], m[1]
		repl := r.to.Interpolate(r.from, word, m)
		rs = append(rs, Replacement{start, stop, repl})
	}

	return rs
}
