package hangulize

import (
	"fmt"
	"strings"
)

// Replacement lets Rewrite know which letters should be replaced.
type Replacement struct {
	start int
	stop  int
	repls []string
}

func (r Replacement) String() string {
	return fmt.Sprintf(`%d-%d -> %s`,
		r.start, r.stop, strings.Join(r.repls, ","),
	)
}

// Replacer determines replacements.
type Replacer interface {
	Replacements(string) []Replacement
}

// -----------------------------------------------------------------------------

type ReplacedBuilder struct {
	flags []bool
}

func (b *ReplacedBuilder) Flags() []bool {
	return b.flags
}

func (b *ReplacedBuilder) Write(flag bool, length int) {
	for i := 0; i < length; i++ {
		b.flags = append(b.flags, flag)
	}
}

func (b *ReplacedBuilder) WriteFrom(flags []bool, start, stop int) {
	for i := start; i < stop; i++ {
		b.flags = append(b.flags, flags[i])
	}
}

func FormatFlags(flags []bool) string {
	var buf strings.Builder

	for _, flag := range flags {
		if flag {
			buf.WriteRune('#')
		} else {
			buf.WriteRune('_')
		}
	}

	return buf.String()
}

// -----------------------------------------------------------------------------

// Rewrite applies multiple replacers on a word.
func Rewrite(word string, reps []Replacer, replaced []bool) (string, []bool) {
	if replaced == nil {
		replaced = make([]bool, len(word))
	} else if len(replaced) != len(word) {
		panic("length of replaced different with word's length")
	}

	for _, rep := range reps {
		var buf strings.Builder
		var rbuf ReplacedBuilder

		offset := 0
		for _, r := range rep.Replacements(word) {
			// TODO(sublee): Support multiple replacements.
			repl := r.repls[0]

			buf.WriteString(word[offset:r.start])
			rbuf.WriteFrom(replaced, offset, r.start)

			buf.WriteString(repl)
			rbuf.Write(true, len(repl))

			offset = r.stop
		}
		buf.WriteString(word[offset:])
		rbuf.WriteFrom(replaced, offset, len(word))

		word = buf.String()
		replaced = rbuf.Flags()
	}

	return word, replaced
}

// -----------------------------------------------------------------------------

// Rule is a replacer based on Pattern and RPatterns.
type Rule struct {
	from *Pattern
	to   []*RPattern
}

// NewRule creates a Rule.
func NewRule(from *Pattern, to ...*RPattern) *Rule {
	return &Rule{from, to}
}

// Replacements determines which letters should be replaced
// based on Pattern and RPatterns.
func (r *Rule) Replacements(word string) []Replacement {
	rs := make([]Replacement, 0)

	for _, m := range r.from.Find(word, -1) {
		start, stop := m[0], m[1]

		repls := make([]string, len(r.to))
		for i, rp := range r.to {
			repls[i] = rp.Interpolate(r.from, word, m)
		}

		rs = append(rs, Replacement{start, stop, repls})
	}

	return rs
}
