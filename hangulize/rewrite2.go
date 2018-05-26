package hangulize

import "strings"

type Replacement struct {
	start int
	stop  int
	words []string
}

type Replacer interface {
	Replacements(string) []Replacement
}

func Rewrite(word string, reps ...Replacer) []string {
	for _, rep := range reps {
		var buf strings.Builder

		offset := 0
		for _, r := range rep.Replacements(word) {
			buf.WriteString(word[offset:r.start])
			buf.WriteString(r.words[0])
			offset = r.stop
		}
		buf.WriteString(word[offset:])

		word = buf.String()
	}

	return []string{word}
}

// -----------------------------------------------------------------------------

type rule2 struct {
	from *Pattern
	to   []*RPattern
}

func newRule(from *Pattern, to ...*RPattern) *rule2 {
	return &rule2{from, to}
}

// Rewrite rewrites a word for a rule.
func (r *rule2) Replacements(word string) []Replacement {
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
