package hsl

import "errors"

var errIllegalToken = errors.New("illegal token")

// token represents a meaningful string in HSL format.
type token int

const (
	// Illegal represents any string not matched with legal tokens.
	Illegal token = iota

	// EOF represents the end-of-file.
	EOF

	// Space represents any of whitespace characters except "\n".
	Space

	// Comment represents a comment content excluding initial "#".
	Comment

	// Newline means only "\n". HSL is a line-sensitive format.
	Newline

	// String represents a text literal.
	String

	// Colon means only ":".
	Colon

	// Comma means only ",".
	Comma

	// Equal means only "=".
	Equal

	// Arrow means only "->".
	Arrow
)
