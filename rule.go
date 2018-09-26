package hangulize

import (
	"fmt"

	"github.com/hangulize/hre"
)

// Rule is a pair of Pattern and RPattern.
type Rule struct {
	ID   int
	From *hre.Pattern
	To   *hre.RPattern
}

func (r Rule) String() string {
	return fmt.Sprintf("%s -> %s", r.From, r.To)
}

// replacements indicates which ranges should be replaced.
func (r Rule) replacements(word string) []Replacement {
	var repls []Replacement

	for _, m := range r.From.Find(word, -1) {
		start, stop := m[0], m[1]

		repl, err := r.To.Interpolate(r.From, word, m)

		if err != nil {
			// FIXME(sublee): Shouldn't it throw an error?
			continue
		}

		repls = append(repls, Replacement{start, stop, repl})
	}

	return repls
}

// Replace matches the word with the Pattern and replaces with the RPattern.
func (r Rule) Replace(word string) string {
	rep := NewReplacer(word, 0, 0)
	repls := r.replacements(word)
	rep.ReplaceBy(repls...)
	return rep.String()
}
