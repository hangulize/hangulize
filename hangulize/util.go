package hangulize

import (
	"regexp"
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

func safeSlice(s string, start int, stop int) string {
	if start < 0 || stop < 0 {
		return ""
	}
	if stop-start > 0 {
		return s[start:stop]
	}
	return ""
}
