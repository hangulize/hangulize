package hgl

import (
	"fmt"
)

// token represents a meaningful string in HGL format.
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

	// Newline means only "\n". HGL is a line-sensitive format.
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

var tokenNames = map[token]string{
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

// formatTokenLiteral formats return value (token, literal) from Scan() as a
// human-readable string.
func formatTokenLiteral(tok token, lit string) string {
	tokenName := tokenNames[tok]
	return fmt.Sprintf(`<%s: %#v>`, tokenName, lit)
}

// illegalError makes an error for an illegal literal.
func illegalError(lit string) error {
	return fmt.Errorf("unexpected token illegal: %#v", lit)
}
