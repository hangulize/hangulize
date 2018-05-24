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

// -----------------------------------------------------------------------------

type rToken int

const (
	plain rToken = iota
	toVar
	// edge
)

type rPart struct {
	tok rToken
	lit string

	// References to the var.
	usedVar []string
}

// -----------------------------------------------------------------------------

// NewRPattern parses the given expression and creates an RPattern.
func NewRPattern(expr string,

	macros map[string]string,
	vars map[string][]string,

) *RPattern {

	_expr := expandMacros(expr, macros)

	// Split expr into several parts.
	// Adjoining 2 parts have different token with each other.
	offset := 0
	parts := make([]rPart, 0)

	for _, m := range reVar.FindAllStringSubmatchIndex(_expr, -1) {
		// Keep plain text before var.
		plainText := _expr[offset:m[0]]
		if plainText != "" {
			parts = append(parts, rPart{plain, plainText, nil})
		}

		// Keep var and the var values.
		varExpr := captured(_expr, m, 0)
		_, vals := getVar(varExpr, vars)
		parts = append(parts, rPart{toVar, varExpr, vals})

		offset = m[1]
	}

	// Keep remaining plain text.
	plainText := _expr[offset:]
	if plainText != "" {
		parts = append(parts, rPart{plain, plainText, nil})
	}

	return &RPattern{expr, parts}
}
