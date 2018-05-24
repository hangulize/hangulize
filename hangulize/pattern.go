package hangulize

import (
	"fmt"
	"regexp"
	"strings"
)

// Pattern represents an HRE (Hangulize-specific Regular Expression) pattern.
// It is used for the rewrite of Hangulize.
//
// Some expressions in Pattern have special meaning:
//
// - "^" - start of chunk
// - "^^" - start of string
// - "$" - end of chunk
// - "$$" - end of string
// - "{...}" - zero-width match
// - "{~...}" - zero-width negative match
// - "<var>" - one of var values (defined in spec)
//
type Pattern struct {
	expr string

	re  *regexp.Regexp // positive regexp
	neg *regexp.Regexp // negative regexp
}

func (p *Pattern) String() string {
	return fmt.Sprintf(`/%s/`, p.expr)
}

// ExplainPattern shows the HRE expression with
// the underlying standard regexp patterns.
func ExplainPattern(p *Pattern) string {
	if p == nil {
		return fmt.Sprintf("%#v", nil)
	}
	return fmt.Sprintf("expr:/%s/, re:/%s/, neg:/%s/", p.expr, p.re, p.neg)
}

// CompilePattern compiles an HRE pattern from an expression.
func CompilePattern(
	expr string,

	macros map[string]string,
	vars map[string][]string,

) (*Pattern, error) {

	reExpr := expr

	// macros
	reExpr = expandMacros(reExpr, macros)

	// vars
	reExpr = expandVars(reExpr, vars)

	// lookaround
	reExpr, negExpr := expandLookbehind(reExpr)
	reExpr, negExpr = expandLookahead(reExpr, negExpr)

	err := mustNoZeroWidth(reExpr)
	if err != nil {
		return nil, err
	}

	if negExpr == `` {
		// This regexp has a paradox.  So it never matches with any text.
		negExpr = `$^`
	}

	// edges
	reExpr = expandEdges(reExpr)

	// Compile regexp.
	re := regexp.MustCompile(reExpr)
	neg := regexp.MustCompile(negExpr)

	p := &Pattern{expr, re, neg}
	return p, nil
}

// expandMacros replaces macro sources to corresponding targets.
// It must be evaluated at the first in CompilePattern.
func expandMacros(reExpr string, macros map[string]string) string {
	args := make([]string, len(macros)*2)

	i := 0
	for src, dst := range macros {
		args[i] = src
		i++
		args[i] = dst
		i++
	}

	replacer := strings.NewReplacer(args...)
	return replacer.Replace(reExpr)
}

// expandVars replaces <var> to corresponding content regexp such as (a|b|c).
func expandVars(reExpr string, vars map[string][]string) string {
	return reVar.ReplaceAllStringFunc(reExpr, func(matched string) string {
		// Retrieve variable name and values.
		name, vals := getVar(matched, vars)

		// Build as RegExp like /(a|b|c)/
		escapedVals := make([]string, len(vals))
		for i, val := range vals {
			escapedVals[i] = regexp.QuoteMeta(val)
		}
		// return `(?P<` + name + `>` + strings.Join(escapedVals, `|`) + `)`
		_ = name
		return `(?:` + strings.Join(escapedVals, `|`) + `)`
	})
}

// Lookbehind: {...} on the left-side.
func expandLookbehind(reExpr string) (string, string) {
	// ^{han}gul
	// │  │   └─ other
	// │  └─ look
	// └─ edge

	posExpr := reExpr
	negExpr := ``

	// This pattern always matches.
	//  [start, stop, edgeStart, edgeStop, lookStart, lookStop]
	m := reLookbehind.FindStringSubmatchIndex(posExpr)

	edgeExpr := safeSlice(posExpr, m[2], m[3])
	lookExpr := safeSlice(posExpr, m[4], m[5])
	otherExpr := posExpr[m[1]:]

	if strings.HasPrefix(lookExpr, `~`) {
		// negative lookbehind
		negExpr = fmt.Sprintf(`(%s)%s`, lookExpr[1:], otherExpr)

		lookExpr = `.*` // require greedy matching
	}

	// Replace lookbehind with 2 parentheses:
	//  (^)(han)gul
	posExpr = fmt.Sprintf(`(%s)(%s)%s`, edgeExpr, lookExpr, otherExpr)

	return posExpr, negExpr
}

// Lookahead: {...} on the right-side.
// negExpr should be passed from expandLookbehind.
func expandLookahead(reExpr string, negExpr string) (string, string) {
	// han{gul}$
	//  │   │  └─ edge
	//  │   └─ look
	//  └─ other

	posExpr := reExpr

	// This pattern always matches:
	//  [start, stop, edgeStart, edgeStop, lookStart, lookStop]
	m := reLookahead.FindStringSubmatchIndex(posExpr)

	otherExpr := posExpr[:m[0]]
	lookExpr := safeSlice(posExpr, m[2], m[3])
	edgeExpr := safeSlice(posExpr, m[4], m[5])

	// Lookahead can be remaining in the negative regexp
	// the lookbehind determined.
	if negExpr != `` {
		lookaheadLen := len(posExpr) - m[0]
		negExpr = negExpr[:len(negExpr)-lookaheadLen]
	}

	if strings.HasPrefix(lookExpr, `~`) {
		// negative lookahead
		if negExpr != `` {
			negExpr += `|`
		}
		negExpr += fmt.Sprintf(`^%s(%s)`, otherExpr, lookExpr[1:])

		lookExpr = `.*` // require greedy matching
	}

	// Replace lookahead with 2 parentheses:
	//  han(gul)($)
	posExpr = fmt.Sprintf(`%s(%s)(%s)`, otherExpr, lookExpr, edgeExpr)

	return posExpr, negExpr
}

func mustNoZeroWidth(reExpr string) error {
	if reZeroWidth.MatchString(reExpr) {
		return fmt.Errorf("zero-width group found in middle: %#v", reExpr)
	}
	return nil
}

func expandEdges(reExpr string) string {
	reExpr = reLeftEdge.ReplaceAllStringFunc(reExpr, func(e string) string {
		switch e {
		case ``:
			return ``
		case `^`:
			return `(?:^|\s+)`
		default:
			// ^^...
			return `^`
		}
	})
	reExpr = reRightEdge.ReplaceAllStringFunc(reExpr, func(e string) string {
		switch e {
		case ``:
			return ``
		case `$`:
			return `(?:$|\s+)`
		default:
			// $$...
			return `$`
		}
	})
	return reExpr
}

// Find searches up to n matches in the word.
func (p *Pattern) Find(word string, n int) [][]int {
	matches := make([][]int, 0)
	offset := 0

	for n < 0 || len(matches) < n {
		// Erase visited characters on the word with "\x00".  Because of
		// lookaround, the search cursor should be calculated manually.
		erased := strings.Repeat(".", offset) + word[offset:]

		m := p.re.FindStringSubmatchIndex(erased)
		if len(m) == 0 {
			// No more match.
			break
		}

		// p.re looks like (edge)(look)abc(look)(edge).
		// Hold only non-zero-width matches.
		start := m[5]
		if start == -1 {
			start = m[0]
		}
		stop := m[len(m)-4]
		if stop == -1 {
			stop = m[1]
		}

		// Pick matched word.  Call it "highlight".
		highlight := erased[m[0]:m[1]]

		// Test highlight with p.neg to determine whether skip or not.
		negM := p.neg.FindStringSubmatchIndex(highlight)

		// If no negative match, this match is successful.
		if len(negM) == 0 {
			matches = append(matches, []int{start, stop})
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
