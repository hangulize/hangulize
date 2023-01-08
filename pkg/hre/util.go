package hre

import (
	"regexp"
	"strings"
)

// indexOf finds the index of the given value in a string array. It returns -1
// if not found. The time complexity is O(n).
func indexOf(val string, vals []string) int {
	for i, _val := range vals {
		if _val == val {
			return i
		}
	}
	return -1
}

// substr is a safe version of s[start:stop]. When start or stop is invalid,
// this function returns "" instead of panic().
func substr(s string, start int, stop int) string {
	if start < 0 || stop < 0 {
		return ""
	}

	n := len(s)
	if start >= n {
		return ""
	}
	if stop > n {
		stop = n
	}

	if stop-start > 0 {
		return s[start:stop]
	}
	return ""
}

// captured returns the captured substring by their group number.
func captured(s string, m []int, n int) string {
	i := n * 2
	return substr(s, m[i], m[i+1])
}

// noCapture removes capturing groups in a regexp string.
func noCapture(expr string) string {
	return strings.Replace(expr, "(", "(?:", -1)
}

// -----------------------------------------------------------------------------
// Verbose Regexp

var (
	// Match with a line starting with "---".
	reComment = regexp.MustCompile(`---.*`)

	// Match with a character or start of string before a whitespace.
	reWhitespace = regexp.MustCompile(`(^|[^\\])\s+`)
)

// re compiles a verbose regular expression.
//
// The expression can be indented and described by comments. Every comment
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

// -----------------------------------------------------------------------------
// Finding meaningful letters from Regexp

var (
	reSpace  = regexp.MustCompile(`\s`)
	reGroup  = regexp.MustCompile(`\(\?(:|P<.+?>)`)
	reMeta   = regexp.MustCompile(`/`)
	reQuoted = regexp.MustCompile(`\\.`)
)

func regexpLetters(reExpr string) string {
	letters := reExpr

	// Remove spaces.
	letters = reSpace.ReplaceAllString(letters, ``)

	// Remove group starters "(?:", "(?P<...>".
	letters = reGroup.ReplaceAllString(letters, ``)

	// Remove meta characters.
	letters = reMeta.ReplaceAllString(letters, ``)

	// Remove escaped letters.
	letters = reQuoted.ReplaceAllString(letters, ``)

	// Quote Regexp meta letters.
	letters = regexp.QuoteMeta(letters)

	// Remove escaped letters again.
	letters = reQuoted.ReplaceAllString(letters, ``)

	return letters
}

func splitLetters(word string) []string {
	var letters []string
	for _, ch := range word {
		letters = append(letters, string(ch))
	}
	return letters
}
