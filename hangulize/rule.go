package hangulize

import "fmt"

// Rule is a pair of Pattern and RPattern.
type Rule struct {
	from *Pattern
	to   *RPattern
}

func (r *Rule) String() string {
	return fmt.Sprintf("%s -> %s", r.from, r.to)
}

// replacements indicates which ranges should be replaced.
func (r *Rule) replacements(word string) []replacement {
	rs := make([]replacement, 0)

	for _, m := range r.from.Find(word, -1) {
		start, stop := m[0], m[1]
		repl := r.to.Interpolate(r.from, word, m)
		rs = append(rs, replacement{start, stop, repl})
	}

	return rs
}
