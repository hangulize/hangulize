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
			to:   make([]string, 0),
		}
	}

	return &Rewriter{rules}, nil
}

func (r *Rewriter) Rewrite(word string) string {
	for i, rule := range r.rules {
		fmt.Println(i, rule)
	}
	return word
}
