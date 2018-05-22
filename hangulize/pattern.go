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
		zeroWidth = `\{(.+?)\}` // {...}
		leftEdge  = `^(\^*)`    // "^", "^^", or empty start
		rightEdge = `(\$*)$`    // "$", "$$", or empty end
	)
	reLookbehind = regexp.MustCompile(leftEdge + zeroWidth)
	reLookahead = regexp.MustCompile(zeroWidth + rightEdge)
	reZeroWidth = regexp.MustCompile(zeroWidth)
}

func expandMacros(reExpr string, spec *Spec) string {
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

func expandVars(reExpr string, spec *Spec) string {
	// Expand <var> through spec.
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

	posExpr := reExpr
	negExpr := ""

	// TODO(sublee): negative lookaround
	// TODO(sublee): edge specialization

	loc = reLookbehind.FindStringSubmatchIndex(posExpr)
	if len(loc) == 6 {
		// lookbehind found

		// ^{han}gul
		// │  │   └─ other
		// │  └─ look
		// └─ edge
		edgeExpr := posExpr[loc[2]:loc[3]]
		lookExpr := posExpr[loc[4]:loc[5]]
		otherExpr := posExpr[loc[1]:]

		// Replace lookbehind with 2 parentheses:
		// (^)(han)gul
		posExpr = fmt.Sprintf(`(%s)(%s)%s`, edgeExpr, lookExpr, otherExpr)
	} else {
		// Prepend empty 2 parentheses.
		posExpr = "()()" + posExpr
	}

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

	if len(negExpr) == 0 {
		negExpr = "$^" // never match
	}

	return posExpr, negExpr
}

// func CompileRewrite(reExpr string, replace func(string) string) Rewrite {
// }
