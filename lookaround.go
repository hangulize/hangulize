package hangulize

import (
	"fmt"
	"strings"
)

var (
	// ^^
	reLeftEdge = re(`\^+`)

	// $$
	reRightEdge = re(`\$+`)

	// {...}
	//  └─┴─ (1)
	reZeroWidth = re(`
	--- open brace
		\{

	--- inside of brace
		( [^}]+ )

	--- close brace
		\}
	`)

	// {...}$$
	//  │ │ └┴─ (2)
	//  └─┴─ (1)
	reLookahead = re(`
	--- zero-width
		(?:
			\{
			( [^}]+ )
			\}
		)?

	--- right-edge
		( \$* )

	--- end of string
		$
	`)

	// ^^{...}
	// ││ └─┴─ (2)
	// └┴─ (1)
	reLookbehind = re(`
	--- start of string
		^

	--- left-edge
		( \^* )

	--- zero-width
		(?:
			\{
			( [^}]+ )
			\}
		)?
	`)
)

func expandLookaround(expr string) (string, string, string, error) {
	posExpr, negAExpr := expandLookahead(expr)
	posExpr, negBExpr := expandLookbehind(posExpr)

	err := mustNoZeroWidth(posExpr)
	if err != nil {
		return ``, ``, ``, err
	}

	return posExpr, negAExpr, negBExpr, nil
}

// Lookahead: {...} on the right-side.
// negExpr should be passed from expandLookbehind.
func expandLookahead(expr string) (string, string) {
	// han{gul}$
	//  │   │  └─ edge
	//  │   └─ look
	//  └─ other

	posExpr := expr
	negAExpr := ``

	// This pattern always matches.
	m := reLookahead.FindStringSubmatchIndex(posExpr)

	start := m[0]
	otherExpr := posExpr[:start]

	edgeExpr := captured(posExpr, m, 2)
	lookExpr := captured(posExpr, m, 1)

	// Don't allow capturing groups in zero-width matches.
	edgeExpr = noCapture(edgeExpr)
	lookExpr = noCapture(lookExpr)

	if strings.HasPrefix(lookExpr, `~`) {
		// negative lookahead
		negAExpr = fmt.Sprintf(`^(%s)`, lookExpr[1:])

		// Lookahead requires greedy matching, unlike lookbehind.
		lookExpr = `.*`
	}

	// Replace lookahead with 2 parentheses:
	//  han(gul)($)
	posExpr = fmt.Sprintf(`%s(%s)(%s)`, otherExpr, lookExpr, edgeExpr)

	return posExpr, negAExpr
}

// Lookbehind: {...} on the left-side.
func expandLookbehind(expr string) (string, string) {
	// ^{han}gul
	// │  │   └─ other
	// │  └─ look
	// └─ edge

	posExpr := expr
	negBExpr := ``

	// This pattern always matches.
	m := reLookbehind.FindStringSubmatchIndex(posExpr)

	stop := m[1]
	otherExpr := posExpr[stop:]

	edgeExpr := captured(posExpr, m, 1)
	lookExpr := captured(posExpr, m, 2)

	// Don't allow capturing groups in zero-width matches.
	edgeExpr = noCapture(edgeExpr)
	lookExpr = noCapture(lookExpr)

	if strings.HasPrefix(lookExpr, `~`) {
		// negative lookbehind
		negBExpr = fmt.Sprintf(`(%s)$`, lookExpr[1:])

		// Lookbehind requires non-greedy matching, unlike lookahead.
		lookExpr = `.*?`
	}

	// Replace lookbehind with 2 parentheses:
	//  (^)(han)gul
	posExpr = fmt.Sprintf(`(%s)(%s)%s`, edgeExpr, lookExpr, otherExpr)

	return posExpr, negBExpr
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
			// "{}" is a zero-width space which is injected by an RPattern.
			return `(?:^|\s+|{})`
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
			// "{}" is a zero-width space which is injected by an RPattern.
			return `(?:$|\s+|{})`
		default:
			// $$...
			return `$`
		}
	})
	return expr
}
