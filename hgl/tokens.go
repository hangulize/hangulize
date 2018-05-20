package hgl

import (
	"fmt"
)

// Token represents a meaningful string in HGL format.
type Token int

const (
	// Illegal represents any string not matched with legal tokens.
	Illegal Token = iota

	// EOF represents the end-of-file.
	EOF

	// Space represents any of whitespace characters except "\n".
	Space

	// Comment represents a comment content excluding initial "#".
	Comment

	// Newline means only "\n".  HGL is a line-sensitive format.
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

var tokenNames = map[Token]string{
	Illegal: `Illegal`,
	EOF:     `EOF`,
	Space:   `Space`,
	Comment: `Comment`,
	Newline: `Newline`,
	String:  `String`,
	Colon:   `Colon`,
	Comma:   `Comma`,
	Equal:   `Equal`,
	Arrow:   `Arrow`,
}

// FormatTokenLiteral formats return value (token, literal) from Scan() as a
// human-readable string.
func FormatTokenLiteral(token Token, literal string) string {
	tokenName := tokenNames[token]
	return fmt.Sprintf(`<%s: %#v>`, tokenName, literal)
}

// IllegalError makes an error for an illegal literal.
func IllegalError(literal string) error {
	return fmt.Errorf("unexpected token illegal: %#v", literal)
}
