package hangulize

import (
	"regexp"
	"sort"
	"strings"
)

var (
	// Match with a line starting with "---".
	reComment = regexp.MustCompile(`---.*`)

	// Match with a character or start of string before a whitespace.
	reWhitespace = regexp.MustCompile(`(^|[^\\])\s+`)
)

// re compiles a verbose regular expression.
//
// The expression can be indented and described by comments.  Every comment
// lines and whitespace except escaped "\ " will be removed before compiling.
//
// Example:
//  var reEmail = re(`
//  --- start of string
//      ^
//  --- user
//      (
//          [^@]+
//      )
//  --- at
//      @
//  --- host
//      (
//          [a-zA-Z0-9-_]+
//          \.
//          [a-zA-Z0-9-_.]+
//      )
//  --- end of string
//      $
//  `)
//
func re(verboseExpr string) *regexp.Regexp {
	expr := reComment.ReplaceAllString(verboseExpr, ``)

	// Remove all whitespace except "\ ".
	expr = reWhitespace.ReplaceAllString(expr, `$1`)

	return regexp.MustCompile(expr)
}

// safeSlice is a safe version of s[start:stop].  When start or stop is
// invalid, this function returns "" instead of panic().
func safeSlice(s string, start int, stop int) string {
	if start < 0 || stop < 0 {
		return ""
	}
	if stop-start > 0 {
		return s[start:stop]
	}
	return ""
}

// captured returns the captured substring by their group number.
func captured(s string, m []int, n int) string {
	i := n * 2
	return safeSlice(s, m[i], m[i+1])
}

// noCapture removes capturing groups in a regexp string.
func noCapture(expr string) string {
	return strings.Replace(expr, "(", "(?:", -1)
}

// indexOf finds the index of the given value in a string array.  It returns -1
// if not found.  The time complexity is O(n).
func indexOf(val string, vals []string) int {
	for i, _val := range vals {
		if _val == val {
			return i
		}
	}
	return -1
}

func set(vals []string) []string {
	unique := make(map[string]bool)
	for _, val := range vals {
		unique[val] = true
	}

	set := make([]string, 0)
	for val := range unique {
		set = append(set, val)
	}

	sort.Strings(set)
	return set
}

func inSet(val string, set []string) bool {
	i := sort.SearchStrings(set, val)
	return i < len(set) && set[i] == val
}

// trueFlags generates [true...] for the given length.
func trueFlags(length int) []bool {
	flags := make([]bool, length)
	for i := 0; i < length; i++ {
		flags[i] = true
	}
	return flags
}
