package hangulize

import (
	"fmt"

	"github.com/hangulize/hre"

	"github.com/hangulize/hangulize/internal/subword"
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

// Replace matches the word with the Pattern and replaces with the RPattern.
func (r Rule) Replace(word string) string {
	rep := subword.NewReplacer(word, 0, 0)
	repls := r.replacements(word)
	rep.ReplaceBy(repls...)
	return rep.String()
}

// replacements indicates which ranges should be replaced.
func (r Rule) replacements(word string) []subword.Replacement {
	var repls []subword.Replacement

	for _, m := range r.From.Find(word, -1) {
		start, stop := m[0], m[1]

		repl, err := r.To.Interpolate(r.From, word, m)

		if err != nil {
			// FIXME(sublee): Shouldn't it throw an error?
			continue
		}

		repls = append(repls, subword.NewReplacement(start, stop, repl))
	}

	return repls
}
