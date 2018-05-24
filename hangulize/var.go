package hangulize

import (
	"regexp"
	"strings"
)

// match: [start, stop]
var reVar = regexp.MustCompile(`<.+?>`)

// expandVars replaces <var> to corresponding content regexp such as (a|b|c).
func expandVars(expr string, vars map[string][]string) string {
	return reVar.ReplaceAllStringFunc(expr, func(matched string) string {
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

func getVar(matched string, vars map[string][]string) (string, []string) {
	// matched looks like "<var>".
	name := strings.Trim(matched, `<>`)
	vals := vars[name]

	return name, vals
}
