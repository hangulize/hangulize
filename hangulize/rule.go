package hangulize

// Rule is a pair of Pattern and RPattern.
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
