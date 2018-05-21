package hangulize

import (
	"regexp"
	"strings"
)

// Pattern is a domain-specific regular expression dialect for Hangulize.
type Pattern struct {
	expr   string
	reExpr string
	re     regexp.Regexp
}

// CompilePattern compiles an Pattern pattern for the given language spec.
func CompilePattern(expr string, spec *Spec) *Pattern {
	reExpr := expr

	reExpr = expandMacros(reExpr, spec)
	reExpr = expandVars(reExpr, spec)

	return &Pattern{expr, reExpr, *regexp.MustCompile(reExpr)}
}

var (
	reVar *regexp.Regexp
)

func init() {
	reVar = regexp.MustCompile("<.+?>")
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

// func CompileRewrite(reExpr string, replace func(string) string) Rewrite {
// }
