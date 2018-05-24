package hangulize

import (
	"fmt"
)

// RPattern is used for dynamic replacement.  "R" of RPattern means
// "replacement" or "right-side".
//
// Some expressions in RPattern have special meaning:
//
// - "/" - zero-width edge of chunk
// - "<var>" - ...
//
type RPattern struct {
	expr string

	parts []rPart
}

func (p *RPattern) String() string {
	return fmt.Sprintf(`"%s"`, p.expr)
}

type rToken int

const (
	plain rToken = iota
	toVar
	edge
)

type rPart struct {
	tok rToken
	lit string

	// References to the var.
	usedVar []string
}

func NewRPattern(expr string,

	macros map[string]string,
	vars map[string][]string,

) (*RPattern, error) {
	// TODO(sublee): RPattern should understand "ab<cd>e" as:
	//
	// - "ab" (normal)
	// - "<cd>" (i: 0, var: cd, vals: c, d)
	// - "e" (norhldkq3al)
	//

	_expr := expr

	_expr = expandMacros(_expr, macros)

	parts := make([]rPart, 0)

	offset := 0
	for _, m := range reVar.FindAllStringSubmatchIndex(_expr, -1) {
		plainText := _expr[offset:m[0]]
		if plainText != "" {
			parts = append(parts, rPart{plain, plainText, nil})
		}

		varExpr := _expr[m[0]:m[1]]
		_, vals := getVar(varExpr, vars)

		parts = append(parts, rPart{toVar, varExpr, vals})

		offset = m[1]
	}

	plainText := _expr[offset:]
	if plainText != "" {
		parts = append(parts, rPart{plain, plainText, nil})
	}

	p := &RPattern{expr, parts}
	return p, nil
}
