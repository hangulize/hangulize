package hangulize

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

// Pattern represents an HRE (Hangulize-specific Regular Expression) pattern.
//
// The transcription logic includes several rewriting rules. A rule has a
// Pattern and an RPattern. A sub-word which is matched with the Pattern, will
// be rewritten by the RPattern.
//
//  rewrite:
//      "'"        -> ""
//      "^gli$"    -> "li"
//      "^glia$"   -> "g.lia"
//      "^glioma$" -> "g.lioma"
//      "^gli{@}"  -> "li"
//      "{@}gli"   -> "li"
//      "gn{@}"    -> "nJ"
//      "gn"       -> "n"
//
// Some expressions in Pattern have special meaning:
//
//  "^"      // start of chunk
//  "^^"     // start of string
//  "$"      // end of chunk
//  "$$"     // end of string
//  "{...}"  // zero-width match
//  "{~...}" // zero-width negative match
//  "{}"     // zero-width space
//  "<var>"  // one of var values (defined in spec)
//
type Pattern struct {
	expr string

	re  *regexp.Regexp // positive regexp
	neg *regexp.Regexp // negative regexp

	// Letters used in the positive/negative regexps.
	letters stringSet

	// References to expanded vars.
	usedVars [][]string
}

func (p *Pattern) String() string {
	return fmt.Sprintf(`/%s/`, p.expr)
}

// newPattern compiles an HRE pattern from an expression.
func newPattern(
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

	reExpr, negExpr, err := expandLookaround(reExpr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compile pattern: %#v", expr)
	}

	reExpr = expandEdges(reExpr)

	// Collect letters in the regexps.
	letters := newStringSet(splitLetters(regexpLetters(reExpr + negExpr))...)

	// Compile regexp.
	re, err := regexp.Compile(reExpr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compile pattern: %#v", expr)
	}

	neg, err := regexp.Compile(negExpr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compile pattern: %#v", expr)
	}

	p := &Pattern{expr, re, neg, letters, usedVars}
	return p, nil
}

// Letters returns the set of natural letters used in the expression in
// ascending order.
func (p *Pattern) Letters() []string {
	return p.letters.Array()
}

// Explain shows the HRE expression with
// the underlying standard regexp patterns.
func (p *Pattern) Explain() string {
	if p == nil {
		return fmt.Sprintf("%#v", nil)
	}
	return fmt.Sprintf("expr:/%s/, re:/%s/, neg:/%s/", p.expr, p.re, p.neg)
}

// -----------------------------------------------------------------------------

// Find searches up to n matches in the word. If n is -1, it will search all
// matches. The result is an array of submatch locations.
func (p *Pattern) Find(word string, n int) [][]int {
	var matches [][]int

	offset := 0
	length := len(word)

	for offset < length && (n < 0 || len(matches) < n) {
		// Erase visited characters on the word with "\x00". Because of
		// lookaround, the search cursor should be calculated manually.
		erased := strings.Repeat("\x00", offset) + word[offset:]

		m := p.re.FindStringSubmatchIndex(erased)
		if len(m) == 0 {
			// No more match.
			break
		}

		// p.re looks like (edge)(look)abc(look)(edge).
		if len(m) < 10 {
			panic(fmt.Errorf("unexpected submatches from %v: %v", p, m))
		}

		// Pick the actual start and stop.
		start, stop := pickStartStop(m)

		// The match MUST NOT be zero-width.
		if stop-start == 0 {
			panic(fmt.Errorf("zero-width match from %v", p))
		}

		// Pick matched word. Call it "highlight".
		highlight := erased[m[0]:m[1]]

		// Test highlight with p.neg to determine whether skip or not.
		negM := p.neg.FindStringSubmatchIndex(highlight)

		// If no negative match, this match is successful.
		if len(negM) == 0 {
			match := []int{start, stop}

			// Keep content ()...
			match = append(match, m[6:len(m)-4]...)

			matches = append(matches, match)
		}

		// Shift the cursor.
		if len(negM) == 0 {
			offset = stop
		} else {
			offset = m[0] + negM[1]
		}
	}

	return matches
}

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
