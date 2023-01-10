package hre

import (
	"fmt"
	"regexp/syntax"
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

func expandLookaround(expr string) (string, string, string, int, int, error) {
	posExpr, negAExpr, negAWidth := expandLookahead(expr)
	posExpr, negBExpr, negBWidth := expandLookbehind(posExpr)

	err := mustNoZeroWidth(posExpr)
	if err != nil {
		return ``, ``, ``, 0, 0, err
	}

	return posExpr, negAExpr, negBExpr, negAWidth, negBWidth, nil
}

func dissolveLookaround(
	format string,
	lookExpr string,
	hasEdge bool,
) (string, string, int, bool) {
	// No Lookaround
	if lookExpr == `` {
		return ``, ``, 0, true
	}

	// Negative Lookaround
	isNeg := strings.HasPrefix(lookExpr, `~`)

	if isNeg {
		var (
			negExpr  string
			negWidth int
		)

		if !hasEdge {
			negExpr = fmt.Sprintf(format, lookExpr[1:])

			re, err := syntax.Parse(negExpr, syntax.Perl)
			if err != nil {
				return ``, ``, 0, false
			}
			negWidth = RegexpMaxWidth(re)
		}

		return ``, negExpr, negWidth, true
	}

	// Positive Lookaround

	if hasEdge {
		// Positive lookaround with edge has a paradox.
		// It wouldn't match on anything.
		return ``, ``, 0, false
	}

	return lookExpr, ``, 0, true
}

// Lookahead: {...} on the right-side.
// negExpr should be passed from expandLookbehind.
func expandLookahead(expr string) (string, string, int) {
	// han{gul}$
	//  │   │  └─ edge
	//  │   └─ look
	//  └─ other

	posExpr := expr

	// This pattern always matches.
	m := reLookahead.FindStringSubmatchIndex(posExpr)

	start := m[0]
	otherExpr := posExpr[:start]

	edgeExpr := captured(posExpr, m, 2)
	lookExpr := captured(posExpr, m, 1)

	// Don't allow capturing groups in zero-width matches.
	edgeExpr = noCapture(edgeExpr)
	lookExpr = noCapture(lookExpr)

	// Dissolve lookahead.
	lookExpr, negAExpr, negAWidth, ok :=
		dissolveLookaround(`^(%s)`, lookExpr, edgeExpr != ``)
	if !ok {
		return `.^^`, ``, 0
	}

	// Replace lookahead with 2 parentheses:
	//  han(gul)($)
	posExpr = fmt.Sprintf(`%s(%s)(%s)`, otherExpr, lookExpr, edgeExpr)

	return posExpr, negAExpr, negAWidth
}

// Lookbehind: {...} on the left-side.
func expandLookbehind(expr string) (string, string, int) {
	// ^{han}gul
	// │  │   └─ other
	// │  └─ look
	// └─ edge

	posExpr := expr

	// This pattern always matches.
	m := reLookbehind.FindStringSubmatchIndex(posExpr)

	stop := m[1]
	otherExpr := posExpr[stop:]

	edgeExpr := captured(posExpr, m, 1)
	lookExpr := captured(posExpr, m, 2)

	// Don't allow capturing groups in zero-width matches.
	edgeExpr = noCapture(edgeExpr)
	lookExpr = noCapture(lookExpr)

	// Dissolve lookbehind.
	lookExpr, negBExpr, negBWidth, ok :=
		dissolveLookaround(`(%s)$`, lookExpr, edgeExpr != ``)
	if !ok {
		return `.^^`, ``, 0
	}

	// Replace lookbehind with 2 parentheses:
	//  (^)(han)gul
	posExpr = fmt.Sprintf(`(%s)(%s)%s`, edgeExpr, lookExpr, otherExpr)

	return posExpr, negBExpr, negBWidth
}

func mustNoZeroWidth(expr string) error {
	if reZeroWidth.MatchString(expr) {
		return fmt.Errorf("zero-width group found in middle: %#v", expr)
	}
	return nil
}

func expandEdges(expr string) string {
	expr = reLeftEdge.ReplaceAllStringFunc(expr, func(e string) string {
		if e == `^` {
			// "{}" is a zero-width space which is injected by an RPattern.
			return `(?:^|\s+|{})`
		}
		// ^^...
		return `^`
	})
	expr = reRightEdge.ReplaceAllStringFunc(expr, func(e string) string {
		if e == `$` {
			// "{}" is a zero-width space which is injected by an RPattern.
			return `(?:$|\s+|{})`
		}
		// $$...
		return `$`
	})
	return expr
}
