package hgl

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

// Lexer reads a bytes buffer and pops the peak token and literal.
type Lexer struct {
	r *bufio.Reader
}

// NewLexer creates a Lexer.
func NewLexer(r io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(r)}
}

const eof = rune(0)

// read reads the rune on the buffer cursor.
func (l *Lexer) read() rune {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread rewinds the buffer cursor once.
func (l *Lexer) unread() {
	_ = l.r.UnreadRune()
}

// readWhile reads runes as a string during
// test function for each rune returns true.
func (l *Lexer) readWhile(test func(rune) bool) string {
	var buf bytes.Buffer

	for ch := l.read(); ch != eof && test(ch); ch = l.read() {
		buf.WriteRune(ch)
	}
	l.unread()

	return buf.String()
}

// space

func isSpace(ch rune) bool {
	return ch != '\n' && unicode.IsSpace(ch)
}

func (l *Lexer) scanSpace() (Token, string) {
	return Space, l.readWhile(isSpace)
}

// comment

func isCommentLetter(ch rune) bool {
	return ch == '#'
}

func isInLine(ch rune) bool {
	return ch != '\n' && ch != eof
}

func (l *Lexer) scanComment() (Token, string) {
	var buf bytes.Buffer
	nEmpty := 0

	for i := 0; ; i++ {
		// Read until visiting a sharp.
		l.readWhile(isSpace)
		ch := l.read()
		if ch != '#' {
			l.unread()
			break
		}

		// Read a line without newline.
		line := l.readWhile(isInLine)
		l.read() // discard remaining newline

		line = strings.TrimSpace(line)

		if line == "" {
			// empty line
			nEmpty++
			continue
		}

		if i > 0 {
			if nEmpty > 0 {
				buf.WriteString("\n\n")
			} else {
				buf.WriteRune(' ')
			}
		}

		nEmpty = 0

		buf.WriteString(line)
	}

	return Comment, buf.String()
}

// string

func isInitialLetter(ch rune) bool {
	return ch == '_' || unicode.IsLetter(ch)
}

func isLetter(ch rune) bool {
	return isInitialLetter(ch) || unicode.IsDigit(ch)
}

func (l *Lexer) scanString() (Token, string) {
	return String, l.readWhile(isLetter)
}

func (l *Lexer) scanQuotedString() (Token, string) {
	var buf bytes.Buffer

	// discard initial quote
	ch := l.read()
	if ch != '"' {
		panic("not a quote")
	}

	escaped := false

	for {
		ch := l.read()

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

func (l *Lexer) scanArrow() (Token, string) {
	first := l.read()
	second := l.read()

	if first != '-' || second != '>' {
		panic("not ->")
	}

	return Arrow, "->"
}

// Scan reads the buffer and returns the peak token and literal.
func (l *Lexer) Scan() (Token, string) {
	ch := l.read()

	if ch == eof {
		return EOF, string(eof)
	}
	if ch == '\n' {
		return Newline, "\n"
	}

	if unicode.IsSpace(ch) {
		l.unread()
		return l.scanSpace()
	}

	if isInitialLetter(ch) {
		l.unread()
		return l.scanString()
	}
	if ch == '"' {
		l.unread()
		return l.scanQuotedString()
	}
	if isCommentLetter(ch) {
		l.unread()
		return l.scanComment()
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
		l.unread()
		return l.scanArrow()
	}

	return Illegal, string(ch)
}
