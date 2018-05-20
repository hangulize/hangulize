package hgl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)
import ()

// Token represents a meaningful string in HGL format.
type Token int

const (
	// Illegal represents any string not matched with legal tokens.
	Illegal Token = iota

	// EOF represents the end-of-file.
	EOF

	// Newline means only "\n".  HGL is a line-sensitive format.
	Newline

	// Space represents any of whitespace characters except "\n".
	Space

	// String represents a text literal.
	String

	// Comment represents a comment literal including initial "#".
	Comment

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
	Newline: `Newline`,
	Space:   `Space`,
	String:  `String`,
	Comment: `Comment`,
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

// Scanner reads a bytes buffer and pops the peak token and literal.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner creates a Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

const eof = rune(0)

// read reads the rune on the buffer cursor.
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread rewinds the buffer cursor once.
func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

// readWhile reads runes as a string during
// test function for each rune returns true.
func (s *Scanner) readWhile(test func(rune) bool) string {
	var buf bytes.Buffer

	for ch := s.read(); test(ch); ch = s.read() {
		buf.WriteRune(ch)
	}
	s.unread()

	return buf.String()
}

// space

func isSpace(ch rune) bool {
	return ch != '\n' && unicode.IsSpace(ch)
}

func (s *Scanner) scanSpace() (Token, string) {
	return Space, s.readWhile(isSpace)
}

// string

func isInitialLetter(ch rune) bool {
	return ch == '_' || unicode.IsLetter(ch)
}

func isLetter(ch rune) bool {
	return isInitialLetter(ch) || unicode.IsDigit(ch)
}

func (s *Scanner) scanString() (Token, string) {
	return String, s.readWhile(isLetter)
}

func (s *Scanner) scanQuotedString() (Token, string) {
	var buf bytes.Buffer

	// discard initial quote
	ch := s.read()
	if ch != '"' {
		panic("not a quote")
	}

	escaped := false

	for {
		ch := s.read()

		if ch == '"' {
			if escaped {
				escaped = false
			} else {
				break
			}
		}

		if ch == '\\' {
			escaped = true
			continue
		}
		if escaped {
			continue
		}

		buf.WriteRune(ch)
	}

	return String, buf.String()
}

// delimiters

func (s *Scanner) scanArrow() (Token, string) {
	first := s.read()
	second := s.read()

	if first != '-' || second != '>' {
		panic("not ->")
	}

	return Arrow, "->"
}

// Scan reads the buffer and returns the peak token and literal.
func (s *Scanner) Scan() (Token, string) {
	ch := s.read()

	if ch == eof {
		return EOF, string(eof)
	}
	if ch == '\n' {
		return Newline, "\n"
	}

	if unicode.IsSpace(ch) {
		s.unread()
		return s.scanSpace()
	}

	if isInitialLetter(ch) {
		s.unread()
		return s.scanString()
	}
	if ch == '"' {
		s.unread()
		return s.scanQuotedString()
	}

	if ch == ':' {
		return Colon, ":"
	}
	if ch == ',' {
		return Comma, ","
	}
	if ch == '=' {
		return Equal, "="
	}
	if ch == '-' {
		s.unread()
		return s.scanArrow()
	}

	return Illegal, string(ch)
}
