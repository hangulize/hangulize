package hre

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"

	"github.com/pkg/errors"
)

// Pattern represents an HRE (Hangulize-specific Regular Expression) pattern.
//
// The transcription logic includes several rewriting rules. A rule has a
// Pattern and an RPattern. A sub-word which is matched with the Pattern, will
// be rewritten by the RPattern.
//
//	rewrite:
//	    "'"        -> ""
//	    "^gli$"    -> "li"
//	    "^glia$"   -> "g.lia"
//	    "^glioma$" -> "g.lioma"
//	    "^gli{@}"  -> "li"
//	    "{@}gli"   -> "li"
//	    "gn{@}"    -> "nJ"
//	    "gn"       -> "n"
//
// Some expressions in Pattern have special meaning:
//
//	"^"      // start of chunk
//	"^^"     // start of string
//	"$"      // end of chunk
//	"$$"     // end of string
//	"{...}"  // zero-width match
//	"{~...}" // zero-width negative match
//	"<var>"  // one of var values (defined in spec)
type Pattern struct {
	expr string

	re   *regexp.Regexp // positive regexp
	negA *regexp.Regexp // negative ahead regexp
	negB *regexp.Regexp // negative behind regexp

	negAWidth int // max width of nagative lookahead
	negBWidth int // max width of nagative lookbehind

	// Letters used in the positive/negative regexps.
	letters map[rune]bool

	// References to expanded vars.
	usedVars [][]string
}

func (p *Pattern) String() string {
	return p.expr
}

// -----------------------------------------------------------------------------

// NewPattern compiles an HRE pattern from an expression.
func NewPattern(
	expr string,

	macros map[string]string,
	vars map[string][]string,

) (*Pattern, error) {
	if len(expr) == 0 {
		return nil, errors.New("empty pattern not allowed")
	}

	reExpr := expr

	reExpr = expandMacros(reExpr, macros)

	reExpr, usedVars := expandVars(reExpr, vars)

	reExpr, negAExpr, negBExpr, negAWidth, negBWidth, err :=
		expandLookaround(reExpr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compile pattern: %#v", expr)
	}

	reExpr = expandEdges(reExpr)

	// Collect letters in the regexps.
	combinedExpr := reExpr + negAExpr + negBExpr
	letters := make(map[rune]bool)
	for _, let := range splitLetters(regexpLetters(combinedExpr)) {
		letters[let] = true
	}

	// Compile regexp.
	re, err := regexp.Compile(reExpr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compile pattern: %#v", expr)
	}

	// Compile negative lookahead/behind regexps.
	var negA *regexp.Regexp
	var negB *regexp.Regexp

	if negAExpr != `` {
		negA, err = regexp.Compile(negAExpr)
		if err != nil {
			err = errors.Wrapf(err, "failed to compile pattern: %#v", expr)
			return nil, err
		}
	}

	if negBExpr != `` {
		negB, err = regexp.Compile(negBExpr)
		if err != nil {
			err = errors.Wrapf(err, "failed to compile pattern: %#v", expr)
			return nil, err
		}
	}

	p := &Pattern{expr, re, negA, negB, negAWidth, negBWidth, letters, usedVars}
	return p, nil
}

// -----------------------------------------------------------------------------

// Letters returns the set of natural letters used in the expression in
// ascending order.
func (p *Pattern) Letters() []rune {
	letters := make([]rune, 0, len(p.letters))
	for let := range p.letters {
		letters = append(letters, let)
	}
	sort.Slice(letters, func(i, j int) bool {
		return letters[i] < letters[j]
	})
	return letters
}

// Explain shows the HRE expression with
// the underlying standard regexp patterns.
func (p *Pattern) Explain() string {
	if p == nil {
		return fmt.Sprintf("%#v", nil)
	}
	return fmt.Sprintf(
		"expr:/%s/, re:/%s/, negA:/%s/, negB:/%s/",
		p.expr, p.re, p.negA, p.negB,
	)
}

// NegativeLookaroundWidths returns the potential widths of negative lookahead
// and negative lookbehind.
//
// -1 means unlimited. An unlimited negative lookround width leads to a
// polynominal time to match. Otherwise, the match consumes only a linear time.
func (p *Pattern) NegativeLookaroundWidths() (negAWidth int, negBWidth int) {
	negAWidth, negBWidth = p.negAWidth, p.negBWidth
	return
}

// -----------------------------------------------------------------------------

func pickStartStop(m []int) (int, int) {
	start := m[5]
	if start == -1 {
		start = m[0]
	}

	stop := m[len(m)-4]
	if stop == -1 {
		stop = m[1]
	}

	return start, stop
}

// Find searches up to n matches in the word. If n is -1, it will search all
// matches. The result is an array of submatch locations.
func (p *Pattern) Find(word string, n int) [][]int {
	var matches [][]int

	offset := 0
	length := len(word)

	for offset < length && (n < 0 || len(matches) < n) {
		// Find submatches on a shifted word.
		m := p.re.FindStringSubmatchIndex(word[offset:])
		for i := range m {
			if m[i] != -1 {
				m[i] += offset
			}
		}

		// Submatches look like:
		//
		// 0      ┌4   ┌5      -2┐  -1┐
		// └(edge)(look)abc(look)(edge)┐
		//  └2   └3      -4┘  -3┘      1
		//
		lenM := len(m)

		if lenM == 0 {
			// No more match.
			break
		}

		if lenM < 10 {
			// Not expected number of groups.
			panic(fmt.Errorf("unexpected submatches from %v: %v", p, m))
		}

		// Pick the actual start and stop.
		start, stop := pickStartStop(m)

		// The match MUST NOT be zero-width.
		if stop-start == 0 {
			panic(fmt.Errorf("zero-width match from %v", p))
		}

		// Shift the cursor for the next iteration.
		offset = stop

		// Test negative lookahead.
		if p.negA != nil {
			var negAStart, negAStop int

			negAStart = m[lenM-4]
			if p.negAWidth == -1 {
				negAStop = len(word)
			} else {
				negAStop = m[lenM-4] + p.negAWidth
			}

			if p.negA.MatchString(substr(word, negAStart, negAStop)) {
				continue
			}
		}

		// Test negative lookbehind.
		if p.negB != nil {
			var negBStart, negBStop int

			negBStop = m[5]
			if p.negBWidth == -1 {
				negBStart = 0
			} else {
				negBStart = m[5] - p.negBWidth
			}

			if p.negB.MatchString(substr(word, negBStart, negBStop)) {
				continue
			}
		}

		// Successfully matched.
		match := []int{start, stop}
		// Keep submatches in the core match.
		match = append(match, m[6:lenM-4]...)

		// Export the match.
		matches = append(matches, match)
	}

	return matches
}

// Replace finds matches and replaces by the given RPattern.
func (p *Pattern) Replace(word string, rpat *RPattern, n int) string {
	var buf bytes.Buffer
	cur := 0

	for _, m := range p.Find(word, n) {
		start, stop := m[0], m[1]
		buf.WriteString(word[cur:start])

		repl, err := rpat.Interpolate(p, word, m)

		if err == nil {
			buf.WriteString(repl)
		} else {
			buf.WriteString(word[start:stop])
		}

		cur = stop
	}
	buf.WriteString(word[cur:])

	return buf.String()
}
