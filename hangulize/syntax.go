package hangulize

import (
	"fmt"
	"regexp"
	"strings"
)

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

func getVar(matched string, vars map[string][]string) (string, []string) {
	// matched looks like "<var>".
	name := strings.Trim(matched, `<>`)
	vals := vars[name]

	return name, vals
}
