package hangulize

import (
	"fmt"

	"github.com/sublee/hangulize2/hgl"
)

// Rewriter is a container of sequential rewriting rules.
type Rewriter struct {
	rules []rewriteRule
}

// NewRewriter creates a Rewriter from HGL pairs which are read from a spec.
func NewRewriter(
	pairs []hgl.Pair,

	macros map[string]string,
	vars map[string][]string,

) (*Rewriter, error) {

	rules := make([]rewriteRule, len(pairs))

	for i, pair := range pairs {
		from, err := NewPattern(pair.Left(), macros, vars)
		if err != nil {
			return nil, err
		}

		right := pair.Right()
		to := make([]*RPattern, len(right))

		for j, expr := range right {
			to[j] = NewRPattern(expr, macros, vars)
		}

		rules[i] = rewriteRule{from, to}
	}

	return &Rewriter{rules}, nil
}

// Rewrite performs rewriting for every rules sequentially.  Each rewriting
// result will be the input for the next rewriting rule.
func (r *Rewriter) Rewrite(word string, ch chan<- Trace) string {
	for _, rule := range r.rules {
		word = rule.rewrite(word, ch)
	}
	return word
}

// -----------------------------------------------------------------------------

// rewriteRule represents a rewriting rule.  It describes how a word should be
// rewritten.
type rewriteRule struct {
	from *Pattern
	to   []*RPattern
}

// rewrite rewrites a word for a rule.
func (r *rewriteRule) rewrite(word string, ch chan<- Trace) string {
	orig := word
	word = r.from.Replace(word, r.to, -1)[0]
	trace(ch, word, orig, fmt.Sprintf("%s->%s", r.from, r.to[0]))
	return word
}
