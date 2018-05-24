package hangulize

import (
	"fmt"
	"strings"

	"github.com/sublee/hangulize2/hgl"
)

// Rule represents a rewriting rule.  It describes how a word should be
// rewritten.
type Rule struct {
	from *Pattern
	to   []*RPattern
}

// Rewrite rewrites a word for a rule.
func (r *Rule) Rewrite(word string) string {
	return r._Rewrite(word, nil)
}

func (r *Rule) _Rewrite(word string, ch chan<- Trace) string {
	orig := word

	var buf strings.Builder
	offset := 0

	for _, m := range r.from.Find(word, -1) {
		start, stop := m[0], m[1]

		buf.WriteString(word[offset:start])

		// TODO(sublee): Support multiple targets.
		repl := r.to[0]

		// Write replacement instead of the match.
		buf.WriteString(repl.expr)

		offset = stop
	}

	buf.WriteString(word[offset:])

	// 	for {
	// 		loc, ok := r.from.Match(word[offset:])
	// 		if !ok {
	// 			buf.WriteString(word[offset:])
	// 			break
	// 		}

	// 		start := loc[0] + offset
	// 		stop := loc[1] + offset

	// 		buf.WriteString(word[offset:start])

	// 		// TODO(sublee): Support multiple targets.
	// 		repl := r.to[0]

	// 		// Write replacement instead of the match.
	// 		buf.WriteString(repl)

	// 		offset = stop
	// 	}

	word = buf.String()
	trace(ch, word, orig, fmt.Sprintf("%s->%#v", r.from, r.to[0]))
	return word
}

// Rewriter is a container of sequential rewriting rules.
type Rewriter struct {
	rules []Rule
}

// NewRewriter creates a Rewriter from HGL pairs which are read from a spec.
func NewRewriter(
	pairs []hgl.Pair,

	macros map[string]string,
	vars map[string][]string,

) (*Rewriter, error) {

	rules := make([]Rule, len(pairs))

	for i, pair := range pairs {
		from, err := CompilePattern(pair.Left(), macros, vars)
		if err != nil {
			return nil, err
		}

		right := pair.Right()
		to := make([]*RPattern, len(right))

		for j, expr := range right {
			p, err := CompileRPattern(expr, macros, vars)
			if err != nil {
				return nil, err
			}
			to[j] = p
		}

		rules[i] = Rule{from, to}
	}

	return &Rewriter{rules}, nil
}

// Rewrite performs rewriting for every rules sequentially.  Each rewriting
// result will be the input for the next rewriting rule.
func (r *Rewriter) Rewrite(word string) string {
	return r._Rewrite(word, nil)
}

func (r *Rewriter) _Rewrite(word string, ch chan<- Trace) string {
	for _, rule := range r.rules {
		word = rule._Rewrite(word, ch)
	}
	return word
}
