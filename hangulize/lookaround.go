package hangulize

import (
	"fmt"
	"strings"
)

var (
	// ^^
	// └┴─ (1)
	reLeftEdge = regex(`\^+`)

	// $$
	// └┴─ (1)
	reRightEdge = regex(`\$+`)

	// {...}
	//  └─┴─ (1)
	reZeroWidth = regex(`
	# open brace
		\{

	# inside of brace
		( [^}]* )

	# close brace
		\}
	`)

	// ^^{...}
	// ││ └─┴─ (2)
	// └┴─ (1)
	reLookbehind = regex(`
	# start of string
		^

	# left-edge
		( \^* )

	# zero-width
		(?:
			\{
			( [^}]* )
			\}
		)?
	`)

	// {...}$$
	//  │ │ └┴─ (2)
	//  └─┴─ (1)
	reLookahead = regex(`
	# zero-width
		(?:
			\{
			( [^}]* )
			\}
		)?

	# right-edge
		( \$* )

	# end of string
		$
	`)
)

func expandLookaround(expr string) (string, string, error) {
	posExpr, negExpr := expandLookbehind(expr)
	posExpr, negExpr = expandLookahead(posExpr, negExpr)

	err := mustNoZeroWidth(posExpr)
	if err != nil {
		return "", "", err
	}

	if negExpr == `` {
		// This regexp has a paradox.  So it never matches with any text.
		negExpr = `$^`
	}

	return posExpr, negExpr, nil
}

// Lookbehind: {...} on the left-side.
func expandLookbehind(expr string) (string, string) {
	// ^{han}gul
	// │  │   └─ other
	// │  └─ look
	// └─ edge

	posExpr := expr
	negExpr := ``

	// This pattern always matches.
	//  [start, stop, edgeStart, edgeStop, lookStart, lookStop]
	m := reLookbehind.FindStringSubmatchIndex(posExpr)

	edgeExpr := safeSlice(posExpr, m[2], m[3])
	lookExpr := safeSlice(posExpr, m[4], m[5])
	otherExpr := posExpr[m[1]:]

	if strings.HasPrefix(lookExpr, `~`) {
		// negative lookbehind
		negExpr = fmt.Sprintf(`(%s)%s`, lookExpr[1:], otherExpr)

		lookExpr = `.*` // require greedy matching
	}

	// Replace lookbehind with 2 parentheses:
	//  (^)(han)gul
	posExpr = fmt.Sprintf(`(%s)(%s)%s`, edgeExpr, lookExpr, otherExpr)

	return posExpr, negExpr
}

// Lookahead: {...} on the right-side.
// negExpr should be passed from expandLookbehind.
func expandLookahead(expr string, negExpr string) (string, string) {
	// han{gul}$
	//  │   │  └─ edge
	//  │   └─ look
	//  └─ other

	posExpr := expr

	// This pattern always matches:
	//  [start, stop, edgeStart, edgeStop, lookStart, lookStop]
	m := reLookahead.FindStringSubmatchIndex(posExpr)

	otherExpr := posExpr[:m[0]]
	lookExpr := safeSlice(posExpr, m[2], m[3])
	edgeExpr := safeSlice(posExpr, m[4], m[5])

	// Lookahead can be remaining in the negative regexp
	// the lookbehind determined.
	if negExpr != `` {
		lookaheadLen := len(posExpr) - m[0]
		negExpr = negExpr[:len(negExpr)-lookaheadLen]
	}

	if strings.HasPrefix(lookExpr, `~`) {
		// negative lookahead
		if negExpr != `` {
			negExpr += `|`
		}
		negExpr += fmt.Sprintf(`^%s(%s)`, otherExpr, lookExpr[1:])

		lookExpr = `.*` // require greedy matching
	}

	// Replace lookahead with 2 parentheses:
	//  han(gul)($)
	posExpr = fmt.Sprintf(`%s(%s)(%s)`, otherExpr, lookExpr, edgeExpr)

	return posExpr, negExpr
}

func mustNoZeroWidth(expr string) error {
	if reZeroWidth.MatchString(expr) {
		return fmt.Errorf("zero-width group found in middle: %#v", expr)
	}
	return nil
}

func expandEdges(expr string) string {
	expr = reLeftEdge.ReplaceAllStringFunc(expr, func(e string) string {
		switch e {
		case ``:
			return ``
		case `^`:
			return `(?:^|\s+)`
		default:
			// ^^...
			return `^`
		}
	})
	expr = reRightEdge.ReplaceAllStringFunc(expr, func(e string) string {
		switch e {
		case ``:
			return ``
		case `$`:
			return `(?:$|\s+)`
		default:
			// $$...
			return `$`
		}
	})
	return expr
}
