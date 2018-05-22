package hangulize

import (
	"fmt"
	"regexp"
	"strings"
)

// Pattern is a domain-specific regular expression dialect for Hangulize.
type Pattern struct {
	expr string

	re  *regexp.Regexp // positive regexp
	neg *regexp.Regexp // negative regexp

	// skip [2]int // [0]: left skip, [1]: right skip
}

func (p *Pattern) String() string {
	return fmt.Sprintf("/%s/", p.expr)
}

// Match reports whether the pattern matches the given text.
func (p *Pattern) Match(text string) []int {
	offset := 0

	for {
		loc := p.re.FindStringSubmatchIndex(text[offset:])

		// Not matched.
		if len(loc) == 0 {
			return make([]int, 0)
		}

		// p.re looks like (edge)(look)abc(look)(edge).
		// Hold only non-zero-width matches.
		start := offset + loc[5]
		stop := offset + loc[len(loc)-4]

		// Pick matched text.  Call it "highlight".
		highlight := text[loc[0]:loc[1]]

		// Test highlight with p.neg to determine whether skip or not.
		negLoc := p.neg.FindStringIndex(highlight)

		// If no negative match, this match has been succeeded.
		if len(negLoc) == 0 {
			return []int{start, stop}
		}

		// Shift the cursor.
		offset = loc[0] + negLoc[1]
	}
}

// CompilePattern compiles an Pattern pattern for the given language spec.
func CompilePattern(expr string, spec *Spec) *Pattern {
	reExpr := expr

	// spec dependent
	reExpr = expandMacros(reExpr, spec)
	reExpr = expandVars(reExpr, spec)

	// lookaround
	reExpr, negExpr := expandLookaround(reExpr)

	re := regexp.MustCompile(reExpr)
	neg := regexp.MustCompile(negExpr)

	return &Pattern{expr, re, neg}
}

var (
	reVar        *regexp.Regexp
	reLookbehind *regexp.Regexp
	reLookahead  *regexp.Regexp
	reZeroWidth  *regexp.Regexp
)

func init() {
	reVar = regexp.MustCompile(`<.+?>`)

	var (
		zeroWidth = `\{(.*?)\}` // {...}
		leftEdge  = `^(\^*)`    // "^", "^^", or empty start
		rightEdge = `(\$*)$`    // "$", "$$", or empty end
	)
	reLookbehind = regexp.MustCompile(leftEdge + zeroWidth)
	reLookahead = regexp.MustCompile(zeroWidth + rightEdge)
	reZeroWidth = regexp.MustCompile(zeroWidth)
}

// expandMacros replaces macro sources in the spec to corresponding
// destinations.  It must be evaluated at the first in CompilePattern.
func expandMacros(reExpr string, spec *Spec) string {
	if spec == nil {
		return reExpr
	}

	args := make([]string, len(spec.Macros)*2)

	i := 0
	for src, dst := range spec.Macros {
		args[i] = src
		i++
		args[i] = dst
		i++
	}

	replacer := strings.NewReplacer(args...)
	return replacer.Replace(reExpr)
}

// expandVars replaces <var> to corresponding content regexp such as (a|b|c).
func expandVars(reExpr string, spec *Spec) string {
	if spec == nil {
		return reVar.ReplaceAllString(reExpr, "()")
	}

	return reVar.ReplaceAllStringFunc(reExpr, func(matched string) string {
		// Retrieve var from spec.
		varName := strings.Trim(matched, "<>")
		varVals := (*spec).Vars[varName]

		// Build as RegExp like /(a|b|c)/
		escapedVals := make([]string, len(varVals))
		for i, val := range varVals {
			escapedVals[i] = regexp.QuoteMeta(val)
		}
		return "(" + strings.Join(escapedVals, "|") + ")"
	})
}

func expandLookaround(reExpr string) (string, string) {
	var loc []int

	// /$^/ never matches with any text.
	negExprs := [2]string{"$^", "$^"}

	posExpr := reExpr

	// TODO(sublee): edge specialization

	// Lookbehind: Find {...} on the left-side.
	loc = reLookbehind.FindStringSubmatchIndex(posExpr)
	if len(loc) == 6 {
		// ^{han}gul
		// │  │   └─ other
		// │  └─ look
		// └─ edge
		edgeExpr := posExpr[loc[2]:loc[3]]
		lookExpr := posExpr[loc[4]:loc[5]]
		otherExpr := posExpr[loc[1]:]

		if strings.HasPrefix(lookExpr, "~") {
			// negative lookbehind
			negExprs[0] = fmt.Sprintf(`(%s)%s`, lookExpr[1:], otherExpr)
			lookExpr = ".*" // require greedy matching
		}

		// Replace lookbehind with 2 parentheses:
		// (^)(han)gul
		posExpr = fmt.Sprintf(`(%s)(%s)%s`, edgeExpr, lookExpr, otherExpr)
	} else {
		// Prepend empty 2 parentheses.
		posExpr = "()()" + posExpr
	}

	// Lookahead: Find {...} on the right-side.
	loc = reLookahead.FindStringSubmatchIndex(posExpr)
	if len(loc) == 6 {
		// lookahead found

		// han{gul}$
		//  │   │  └─ edge
		//  │   └─ look
		//  └─ other
		otherExpr := posExpr[:loc[0]]
		lookExpr := posExpr[loc[2]:loc[3]]
		edgeExpr := posExpr[loc[4]:loc[5]]

		if strings.HasPrefix(lookExpr, "~") {
			// negative lookahead
			negExprs[1] = fmt.Sprintf(`^%s(%s)`, otherExpr, lookExpr[1:])
			lookExpr = ".*" // require greedy matching
		}

		// Replace lookahead with 2 parentheses:
		// han(gul)($)
		posExpr = fmt.Sprintf(`%s(%s)(%s)`, otherExpr, lookExpr, edgeExpr)
	} else {
		// Append empty 2 parentheses.
		posExpr = posExpr + "()()"
	}

	// Find remaining zero-width groups.
	if reZeroWidth.MatchString(posExpr) {
		panic(fmt.Errorf("zero-width group found in middle: %#v", reExpr))
	}

	negExpr := strings.Join(negExprs[:], "|")
	return posExpr, negExpr
}
