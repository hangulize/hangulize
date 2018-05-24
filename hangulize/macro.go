package hangulize

import (
	"strings"
)

// expandMacros replaces macro sources to corresponding targets.
func expandMacros(expr string, macros map[string]string) string {
	args := make([]string, len(macros)*2)

	i := 0
	for src, dst := range macros {
		args[i] = src
		i++
		args[i] = dst
		i++
	}

	replacer := strings.NewReplacer(args...)
	return replacer.Replace(expr)
}
