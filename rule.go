package hangulize

import (
	"fmt"

	"github.com/hangulize/hre"
)

// Rule is a pair of Pattern and RPattern.
type Rule struct {
	From *hre.Pattern
	To   *hre.RPattern
}

func (r *Rule) String() string {
	return fmt.Sprintf("%s -> %s", r.From, r.To)
}

// replacements indicates which ranges should be replaced.
func (r *Rule) replacements(word string) []replacement {
	var repls []replacement

	for _, m := range r.From.Find(word, -1) {
		start, stop := m[0], m[1]

		repl, err := r.To.Interpolate(r.From, word, m)

		if err != nil {
			repl = word[start:stop]
		}

		repls = append(repls, replacement{start, stop, repl})
	}

	return repls
}
