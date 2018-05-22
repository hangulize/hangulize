package hangulize

import (
	"fmt"
	"regexp"
	"strings"
)

// Pattern represents an HRE (Hangulize-specific Regular Expression) pattern.
// It is used for the rewrite of Hangulize.
type Pattern struct {
	expr string

	re  *regexp.Regexp // positive regexp
	neg *regexp.Regexp // negative regexp
}

func (p *Pattern) String() string {
	return fmt.Sprintf("/%s/", p.expr)
}

func safeSlice(s string, start int, stop int) string {
	if start < 0 || stop < 0 {
		return ""
	}
	if stop-start > 0 {
		return s[start:stop]
	}
	return ""
}

// Match reports whether the pattern matches the given word.
func (p *Pattern) Match(word string) ([]int, bool) {
	offset := 0

	for {
		loc := p.re.FindStringSubmatchIndex(word[offset:])

		// Not matched.
		if len(loc) == 0 {
			return make([]int, 0), false
		}

		// p.re looks like (edge)(look)abc(look)(edge).
		// Hold only non-zero-width matches.
		lookbehindStop := loc[5]
		if lookbehindStop == -1 {
			lookbehindStop = loc[0]
		}
		lookaheadStart := loc[len(loc)-4]
		if lookaheadStart == -1 {
			lookaheadStart = loc[1]
		}
		start := offset + lookbehindStop
		stop := offset + lookaheadStart

		// Pick matched word.  Call it "highlight".
		highlight := word[loc[0]:loc[1]]

		// Test highlight with p.neg to determine whether skip or not.
		negLoc := p.neg.FindStringIndex(highlight)

		// If no negative match, this match has been succeeded.
		if len(negLoc) == 0 {
			return []int{start, stop}, true
		}

		// Shift the cursor.
		offset = loc[0] + negLoc[1]
	}
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

// Pre-compiled regexp patterns to compile HRE patterns.
var (
	reVar        *regexp.Regexp
	reLookbehind *regexp.Regexp
	reLookahead  *regexp.Regexp
	reZeroWidth  *regexp.Regexp
	reLeftEdge   *regexp.Regexp
	reRightEdge  *regexp.Regexp
)

func init() {
	reVar = regexp.MustCompile(`<.+?>`)

	var (
		zeroWidth = `\{([^}]*)\}` // {...}
		leftEdge  = `(\^+)`       // `^`, `^^`, `^^^...`
		rightEdge = `(\$+)`       // `$`, `$$`, `$$$...`

		// begin of text - optional leftEdge - optional zeroWidth
		lookbehind = fmt.Sprintf(`^(?:%s)?(?:%s)?`, leftEdge, zeroWidth)
		// optional zeroWidth - optional rightEdge - end of text
		lookahead = fmt.Sprintf(`(?:%s)?(?:%s)?$`, zeroWidth, rightEdge)
	)

	reLookbehind = regexp.MustCompile(lookbehind)
	reLookahead = regexp.MustCompile(lookahead)
	reZeroWidth = regexp.MustCompile(zeroWidth)

	reLeftEdge = regexp.MustCompile(leftEdge)
	reRightEdge = regexp.MustCompile(rightEdge)
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
		// Retrieve var from spec.
		varName := strings.Trim(matched, `<>`)
		varVals := vars[varName]

		// Build as RegExp like /(a|b|c)/
		escapedVals := make([]string, len(varVals))
		for i, val := range varVals {
			escapedVals[i] = regexp.QuoteMeta(val)
		}
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
	loc := reLookbehind.FindStringSubmatchIndex(posExpr)

	edgeExpr := safeSlice(posExpr, loc[2], loc[3])
	lookExpr := safeSlice(posExpr, loc[4], loc[5])
	otherExpr := posExpr[loc[1]:]

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
	loc := reLookahead.FindStringSubmatchIndex(posExpr)

	otherExpr := posExpr[:loc[0]]
	lookExpr := safeSlice(posExpr, loc[2], loc[3])
	edgeExpr := safeSlice(posExpr, loc[4], loc[5])

	// Lookahead can be remaining in the negative regexp
	// the lookbehind determined.
	if negExpr != `` {
		lookaheadLen := len(posExpr) - loc[0]
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
