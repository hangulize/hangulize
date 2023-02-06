package hre

import (
	"bytes"
	"sort"

	"github.com/pkg/errors"
)

// RPattern is a dynamic replacement pattern.
//
// Some expressions in RPattern have special meaning:
//
//	"{}"    // zero-width space
//	"<var>" // ...
//
// "R" in the name means "replacement" or "right-side".
type RPattern struct {
	expr string

	parts []rPart

	// Letters used in the regexp.
	letters map[rune]bool
}

func (rp *RPattern) String() string {
	return rp.expr
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
func NewRPattern(
	expr string,

	macros map[string]string,
	vars map[string][]string,

) *RPattern {
	_expr := expandMacros(expr, macros)

	// Split expr into several parts.
	// Adjoining 2 parts have different token with each other.
	offset := 0
	var parts []rPart

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

	// Collect letters in the regexp.
	letters := make(map[rune]bool)
	for _, let := range splitLetters(regexpLetters(expr)) {
		letters[let] = true
	}

	return &RPattern{expr, parts, letters}
}

// -----------------------------------------------------------------------------

// Interpolate determines the final replacement based on the matched Pattern.
func (rp *RPattern) Interpolate(
	p *Pattern, word string, m []int,
) (string, error) {
	var buf bytes.Buffer
	varIndex := 0

	for _, part := range rp.parts {
		switch part.tok {

		case plain:
			// just plain text
			buf.WriteString(part.lit)

		case toVar:
			// var-to-var: <var> in Pattern to <var> in RPattern.
			if varIndex > len(p.usedVars)-1 {
				err := errors.New("mapped vars have different length")
				return "", err
			}
			fromVar := p.usedVars[varIndex]
			fromVal := captured(word, m, varIndex+1)

			// Find index of the matched character in the var.
			i := indexOf(fromVal, fromVar)
			i = i % len(part.usedVar)

			// Choose a replacement character at the same index.
			toVal := part.usedVar[i]

			buf.WriteString(toVal)
			varIndex++
		}
	}

	return buf.String(), nil
}

// Letters returns the set of natural letters used in the expression in
// ascending order.
func (rp *RPattern) Letters() []rune {
	letters := make([]rune, 0, len(rp.letters))
	for let := range rp.letters {
		letters = append(letters, let)
	}
	sort.Slice(letters, func(i, j int) bool {
		return letters[i] < letters[j]
	})
	return letters
}
