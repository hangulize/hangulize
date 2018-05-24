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

func CompileRPattern(expr string,

	macros map[string]string,
	vars map[string][]string,

) (*RPattern, error) {

	_expr := expr

	_expr = expandMacros(_expr, macros)

	for _, m := range reVar.FindAllStringSubmatchIndex(_expr, -1) {
		fmt.Println(m)
	}

	p := &RPattern{expr}
	return p, nil
}
