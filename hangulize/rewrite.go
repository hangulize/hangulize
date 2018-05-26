package hangulize

import (
	"strings"
)

// Replacement lets Rewrite know which letters should be replaced.
type Replacement struct {
	start int
	stop  int
	repls []string
}

// Replacer determines replacements.
type Replacer interface {
	Replacements(string) []Replacement
}

// -----------------------------------------------------------------------------

// Rewrite applies multiple replacers on a word.
func Rewrite(word string, reps ...Replacer) []string {
	for _, rep := range reps {
		var buf strings.Builder

		offset := 0
		for _, r := range rep.Replacements(word) {
			buf.WriteString(word[offset:r.start])
			buf.WriteString(r.repls[0])
			offset = r.stop
		}
		buf.WriteString(word[offset:])

		word = buf.String()
	}

	return []string{word}
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
