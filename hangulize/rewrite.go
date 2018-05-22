package hangulize

import (
	"strings"

	"github.com/sublee/hangulize2/hgl"
)

// Rule represents a rewriting rule.  It describes how a word should be
// rewritten.
type Rule struct {
	from *Pattern
	to   []string
}

// Rewrite rewrites a word for a rule.
func (r *Rule) Rewrite(word string) string {
	return r._Rewrite(word, nil)
}

func (r *Rule) _Rewrite(word string, ch chan<- Event) string {
	var buf strings.Builder
	offset := 0

	for {
		loc, ok := r.from.Match(word[offset:])
		if !ok {
			buf.WriteString(word[offset:])
			break
		}

		// fmt.Println(r.from, loc)

		start := loc[0] + offset
		stop := loc[1] + offset

		buf.WriteString(word[offset:start])

		// TODO(sublee): Support multiple targets.
		buf.WriteString(r.to[0])

		offset = stop
	}

	word = buf.String()
	// if offset != 0 {
	// 	fmt.Println(word, r.from, r.to)
	// }
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

// Rewrite performs rewriting for every rules sequentially.  Each rewriting
// result will be the input for the next rewriting rule.
func (r *Rewriter) Rewrite(word string) string {
	return r._Rewrite(word, nil)
}

func (r *Rewriter) _Rewrite(word string, ch chan<- Event) string {
	for _, rule := range r.rules {
		word = rule._Rewrite(word, ch)
	}
	return word
}
