package hgl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)
import ()

type token int

const (
	Illegal token = iota
	EOF
	Space
	Newline
	String
	Comment
	Colon
	Comma
	Equal
	Arrow
)

var tokenNames = map[token]string{
	Illegal: `Illegal`,
	EOF:     `EOF`,
	Space:   `Space`,
	Newline: `Newline`,
	String:  `String`,
	Comment: `Comment`,
	Colon:   `Colon`,
	Comma:   `Comma`,
	Equal:   `Equal`,
	Arrow:   `Arrow`,
}

func FormatTokenLiteral(token token, literal string) string {
	tokenName := tokenNames[token]
	return fmt.Sprintf(`<%s: %#v>`, tokenName, literal)
}

type Scanner struct {
	r *bufio.Reader
}

func newScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

const eof = rune(0)

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() {
	_ = s.r.UnreadRune()
}

func (s *Scanner) scanWhile(test func(rune) bool) string {
	var buf bytes.Buffer

	for ch := s.read(); test(ch); ch = s.read() {
		buf.WriteRune(ch)
	}
	s.unread()

	return buf.String()
}

// space

func (s *Scanner) scanSpace() (token, string) {
	return Space, s.scanWhile(unicode.IsSpace)
}

// string

func isInitialLetter(ch rune) bool {
	return ch == '_' || unicode.IsLetter(ch)
}

func isLetter(ch rune) bool {
	return isInitialLetter(ch) || unicode.IsDigit(ch)
}

func (s *Scanner) scanString() (token, string) {
	return String, s.scanWhile(isLetter)
}

func (s *Scanner) scanQuotedString() (token, string) {
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

func (s *Scanner) scanArrow() (token, string) {
	first := s.read()
	second := s.read()

	if first != '-' || second != '>' {
		panic("not ->")
	}

	return Arrow, "->"
}

func (s *Scanner) Scan() (token, string) {
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
