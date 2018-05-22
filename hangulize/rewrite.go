package hangulize

import (
	"fmt"
	"strings"

	"github.com/sublee/hangulize2/hgl"
)

// Rule ...
type Rule struct {
	from *Pattern
	to   []string
}

func (r *Rule) Rewrite(word string) string {
	var buf strings.Builder
	offset := 0

	for {
		loc, ok := r.from.Match(word[offset:])
		if !ok {
			buf.WriteString(word[offset:])
			break
		}

		fmt.Println(r.from, loc)
		start := loc[0] + offset
		stop := loc[1] + offset

		buf.WriteString(word[offset:start])

		// TODO(sublee): Support multiple targets.
		buf.WriteString(r.to[0])

		offset = stop
	}

	word = buf.String()
	if offset != 0 {
		fmt.Println(word, r.from, r.to)
	}
	return word
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
