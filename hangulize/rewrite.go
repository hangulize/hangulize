package hangulize

import (
	"fmt"
	"strings"

	"github.com/sublee/hangulize2/hgl"
)

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
		from, err := NewPattern(pair.Left(), macros, vars)
		if err != nil {
			return nil, err
		}

		right := pair.Right()
		to := make([]*RPattern, len(right))

		for j, expr := range right {
			to[j] = NewRPattern(expr, macros, vars)
		}

		rules[i] = Rule{from, to}
	}

	return &Rewriter{rules}, nil
}

// Rewrite performs rewriting for every rules sequentially.  Each rewriting
// result will be the input for the next rewriting rule.
func (r *Rewriter) Rewrite(word string, ch chan<- Trace) string {
	for _, rule := range r.rules {
		word = rule.Rewrite(word, ch)
	}
	return word
}

// -----------------------------------------------------------------------------

// Rule represents a rewriting rule.  It describes how a word should be
// rewritten.
type Rule struct {
	from *Pattern
	to   []*RPattern
}

// Rewrite rewrites a word for a rule.
func (r *Rule) Rewrite(word string, ch chan<- Trace) string {
	orig := word

	var buf strings.Builder
	offset := 0

	for _, m := range r.from.Find(word, -1) {
		start, stop := m[0], m[1]
		trace(ch, orig, "", fmt.Sprintf("%v", m))

		buf.WriteString(word[offset:start])

		// TODO(sublee): Support multiple targets.
		to := r.to[0]

		// Write replacement instead of the match.
		buf.WriteString(Interpolate(r.from, to, word, m))

		offset = stop
	}

	buf.WriteString(word[offset:])

	word = buf.String()
	trace(ch, word, orig, fmt.Sprintf("%s->%s", r.from, r.to[0]))
	return word
}

// Interpolate determines the final replacement based on Pattern and RPattern.
// TODO(sublee): Move to Pattern.Replace()
func Interpolate(left *Pattern, right *RPattern, word string, m []int) string {
	var buf strings.Builder

	varIndex := 0

	for _, part := range right.parts {
		switch part.tok {
		case plain:
			buf.WriteString(part.lit)
			break
		case toVar:
			fromVarVals := left.usedVars[varIndex]

			fromVarVal := captured(word, m, varIndex+1)

			pos := 0
			for ; fromVarVals[pos] != fromVarVal; pos++ {
			}

			buf.WriteString(part.usedVar[pos])

			varIndex++
			break
		}
	}

	return buf.String()
}
