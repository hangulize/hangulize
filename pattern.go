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

	re   *regexp.Regexp // positive regexp
	negB *regexp.Regexp // negative behind regexp
	negA *regexp.Regexp // negative ahead regexp

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

	reExpr, negBExpr, negAExpr, err := expandLookaround(reExpr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compile pattern: %#v", expr)
	}

	reExpr = expandEdges(reExpr)

	// Collect letters in the regexps.
	combinedExpr := reExpr + negBExpr + negAExpr
	letters := newStringSet(splitLetters(regexpLetters(combinedExpr))...)

	// Compile regexp.
	re, err := regexp.Compile(reExpr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compile pattern: %#v", expr)
	}

	negB, err := regexp.Compile(negBExpr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compile pattern: %#v", expr)
	}

	negA, err := regexp.Compile(negAExpr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compile pattern: %#v", expr)
	}

	p := &Pattern{expr, re, negB, negA, letters, usedVars}
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
	return fmt.Sprintf(
		"expr:/%s/, re:/%s/, negB:/%s/, negA:/%s/",
		p.expr, p.re, p.negB, p.negA,
	)
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

		//                      0                         1
		//                      │                         │
		// Submatches look like (edge)(look)abc(look)(edge).
		//                      │    ││    │   │    ││    │
		//                      2    34    5  -4   -3-2  -1
		m := p.re.FindStringSubmatchIndex(erased)

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
		start, stop := p.pickStartStop(m)

		// The match MUST NOT be zero-width.
		if stop-start == 0 {
			panic(fmt.Errorf("zero-width match from %v", p))
		}

		// Test negative lookaround.
		var (
			behind = safeSlice(word, m[4], m[5])
			ahead  = safeSlice(word, m[6], m[7])
		)

		neg := p.negB.MatchString(behind) || p.negA.MatchString(ahead)

		if !neg {
			// No negative lookaround matches.
			match := []int{start, stop}
			match = append(match, m[6:len(m)-4]...)

			matches = append(matches, match)
		}

		// Shift the cursor.
		offset = stop
	}

	return matches
}

func (Pattern) pickStartStop(m []int) (int, int) {
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
