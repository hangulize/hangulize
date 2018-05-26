package hangulize

import (
	"fmt"
	"strings"
)

// Replacement lets Rewrite know which letters should be replaced.
type Replacement struct {
	start int
	stop  int
	repls []string
}

func (r Replacement) String() string {
	return fmt.Sprintf(`%d-%d -> %s`,
		r.start, r.stop, strings.Join(r.repls, ","),
	)
}

// Replacer determines replacements.
type Replacer interface {
	Replacements(string) []Replacement
}

// -----------------------------------------------------------------------------

// Replaced is an indicator which letters have been replaced.
type Replaced struct {
	flags []bool
}

func NewReplaced(word string) *Replaced {
	flags := make([]bool, len(word))
	return &Replaced{flags}
}

func (r Replaced) String() string {
	var buf strings.Builder

	for _, flag := range r.flags {
		if flag {
			buf.WriteRune('#')
		} else {
			buf.WriteRune('_')
		}
	}

	return buf.String()
}

func (r *Replaced) Mark(start, stop, length int) {
	left := r.flags[:start]
	right := r.flags[stop:]

	mid := make([]bool, length)
	for i := 0; i < length; i++ {
		mid[i] = true
	}

	r.flags = append(left, append(mid, right...)...)
}

// -----------------------------------------------------------------------------

// Rewrite applies multiple replacers on a word.
func Rewrite(word string, reps ...Replacer) (string, *Replaced) {
	replaced := NewReplaced(word)

	for _, rep := range reps {
		var buf strings.Builder

		offset := 0
		for _, r := range rep.Replacements(word) {
			// TODO(sublee): Support multiple replacements.
			repl := r.repls[0]

			replaced.Mark(r.start, r.stop, len(repl))

			buf.WriteString(word[offset:r.start])
			buf.WriteString(repl)
			offset = r.stop
		}
		buf.WriteString(word[offset:])

		word = buf.String()
	}

	return word, replaced
}

// -----------------------------------------------------------------------------

// Rule is a replacer based on Pattern and RPatterns.
type Rule struct {
	from *Pattern
	to   []*RPattern
}

// NewRule creates a Rule.
func NewRule(from *Pattern, to ...*RPattern) *Rule {
	return &Rule{from, to}
}

// Replacements determines which letters should be replaced
// based on Pattern and RPatterns.
func (r *Rule) Replacements(word string) []Replacement {
	rs := make([]Replacement, 0)

	for _, m := range r.from.Find(word, -1) {
		start, stop := m[0], m[1]

		repls := make([]string, len(r.to))
		for i, rp := range r.to {
			repls[i] = rp.Interpolate(r.from, word, m)
		}

		rs = append(rs, Replacement{start, stop, repls})
	}

	return rs
}
