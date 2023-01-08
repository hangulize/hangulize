/*
Package hre provides the regular expression dialect for Hangulize called HRE.
HRE focuses on a very narrow usage.

The HRE syntax is based on RE2. But it tweaks the assertions. For example, in
HRE ^ matches with every beginning of a word, not only the beginning of a
string.

Lookaround is not supported in RE2 because there's no known efficient algorithm
without backtracking to implement it. Anyways, HRE provides a simplified
lookaround. The syntax {...} is for the positive lookaround and {~...} is for
the negative lookaround. The lookaround is restircted to place at the leftmost
or rightmost.

	"foo{bar}"
	"{~bar}foo"

The time complexity of the negative lookbehind is O(n²) while other assertions
can be done in O(n). The negative lookbehind should not be used for very long
string.

HRE also provides macros and variables.

	macros := map[string]string {
		"@": "<vowels>",
	}

	vars := map[string][]string {
		"abc":    []string{"a", "b", "c"},
		"vowels": []string{"a", "e", "i", "o", "u"},
	}

	p, err := NewPattern("<abc>@", macros, vars)
	// The p matches with "ai", "be" or "ci".

*/
package hre
