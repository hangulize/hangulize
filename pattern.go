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
//  "<var>"  // one of var values (defined in spec)
//
type Pattern struct {
	expr string

	re   *regexp.Regexp // positive regexp
	negA *regexp.Regexp // negative ahead regexp
	negB *regexp.Regexp // negative behind regexp

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

	reExpr, negAExpr, negBExpr, err := expandLookaround(reExpr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compile pattern: %#v", expr)
	}

	reExpr = expandEdges(reExpr)

	// Collect letters in the regexps.
	combinedExpr := reExpr + negAExpr + negBExpr
	letters := newStringSet(splitLetters(regexpLetters(combinedExpr))...)

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

	p := &Pattern{expr, re, negA, negB, letters, usedVars}
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
		"expr:/%s/, re:/%s/, negA:/%s/, negB:/%s/",
		p.expr, p.re, p.negA, p.negB,
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

		// Submatches look like:
		//
		// 0      ┌4   ┌5      -2┐  -1┐
		// └(edge)(look)abc(look)(edge)┐
		//  └2   └3      -4┘  -3┘      1
		//
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

		// Shift the cursor for the next iteration.
		offset = stop

		// Test negative lookahead.
		if p.negA != nil {
			ahead := safeSlice(word, m[lenM-4], m[lenM-3])
			if p.negA.MatchString(ahead) {
				continue
			}
		}

		// Test negative lookbehind.
		if p.negB != nil {
			behind := safeSlice(word, m[4], m[5])
			if p.negB.MatchString(behind) {
				continue
			}
		}

		// Successfully matched.
		match := []int{start, stop}
		// Keep submatches in the core match.
		match = append(match, m[6:lenM-4]...)

		matches = append(matches, match)
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
