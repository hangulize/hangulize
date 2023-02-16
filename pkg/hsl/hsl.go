/*
Package hsl implements a parser for the HSL format which is used for
Hangulize.

The HSL format is a simple config language. A HSL config has sections. Each
section contains string-strings pairs. A section can be one of a dictionary or
a pair list:

	# dictionary section
	lang:
		id      = "ita"
		code    = "it", "ita", "ita"
		english = "Italian"
		korean  = "이탈리아어"
		script  = "roman"

	# pair list section
	rewrite:
		"^gli$"    -> "li"
		"^glia$"   -> "g.lia"
		"^glioma$" -> "g.lioma"
		"^gli{@}"  -> "li"
		"{@}gli"   -> "li"

Keys in a dictionary section must be unique, and the section won't keep their
order. While a pair list section works in an inversed way. A pair list
section just keeps described pairs in therir order.

The media type of HSL files: application/vnd.hsl
*/
package hsl

import (
	"io"
)

// Parse parses an HSL formatted text.
func Parse(r io.Reader) (HSL, error) {
	p := newParser(r)
	return p.parse()
}
