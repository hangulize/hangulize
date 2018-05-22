package hangulize

import (
	"fmt"

	"github.com/sublee/hangulize2/hgl"
)

// Rule ...
type Rule struct {
	from *Pattern
	to   []string
}

func (r *Rule) Rewrite(word string) string {
	loc, ok := r.from.Match(word)
	if !ok {
		return word
	}

	start, stop := loc[0], loc[1]

	fmt.Println(word, r.from, loc, r.to)

	return word[:start] + r.to[0] + word[stop:]
}

// Rewriter ...
type Rewriter struct {
	rules []Rule
}

func NewRewriter(
	pairs []hgl.Pair,

	macros map[string]string,
	vars map[string][]string,

) (*Rewriter, error) {

	rules := make([]Rule, len(pairs))

	for i, pair := range pairs {
		p, err := CompilePattern(pair.Left(), macros, vars)
		if err != nil {
			return nil, err
		}

		rules[i] = Rule{
			from: p,
			to:   pair.Right(),
		}
	}

	return &Rewriter{rules}, nil
}

func (r *Rewriter) Rewrite(word string) string {
	for _, rule := range r.rules {
		word = rule.Rewrite(word)
	}
	return word
}
