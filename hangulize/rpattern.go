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
}

func (p *RPattern) String() string {
	return fmt.Sprintf(`"%s"`, p.expr)
}

type rpToken int

const (
	plain rpToken = iota
	toVar
	edge
)

func CompileRPattern(expr string,

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

	for _, m := range reVar.FindAllStringSubmatchIndex(_expr, -1) {
		fmt.Println(m)
	}

	p := &RPattern{expr}
	return p, nil
}
