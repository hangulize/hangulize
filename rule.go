package hangulize

import "fmt"

// Rule is a pair of Pattern and RPattern.
type Rule struct {
	From *Pattern
	To   *RPattern
}

func (r *Rule) String() string {
	return fmt.Sprintf("%s -> %s", r.From, r.To)
}

// replacements indicates which ranges should be replaced.
func (r *Rule) replacements(word string) []replacement {
	var repls []replacement

	for _, m := range r.From.Find(word, -1) {
		start, stop := m[0], m[1]
		repl := r.To.Interpolate(r.From, word, m)
		repls = append(repls, replacement{start, stop, repl})
	}

	return repls
}
