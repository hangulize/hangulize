package hgl

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"unicode"
)

// lexer reads a bytes buffer and pops the peak token and literal.
type lexer struct {
	r *bufio.Reader
}

// newLexer creates a Lexer.
func newLexer(r io.Reader) *lexer {
	return &lexer{r: bufio.NewReader(r)}
}

const eof = rune(0)

// read reads the rune on the buffer cursor.
func (l *lexer) read() rune {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread rewinds the buffer cursor once.
func (l *lexer) unread() {
	_ = l.r.UnreadRune()
}

// readWhile reads runes as a string during
// test function for each rune returns true.
func (l *lexer) readWhile(test func(rune) bool) string {
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

func (l *lexer) scanSpace() (token, string) {
	return Space, l.readWhile(isSpace)
}

// comment

func isCommentLetter(ch rune) bool {
	return ch == '#'
}

func isInLine(ch rune) bool {
	return ch != '\n' && ch != eof
}

func (l *lexer) scanComment() (token, string) {
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

func (l *lexer) scanString() (token, string) {
	return String, l.readWhile(isLetter)
}

func (l *lexer) scanQuotedString() (token, string) {
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

func (l *lexer) scanArrow() (token, string) {
	first := l.read()
	second := l.read()

	if first != '-' || second != '>' {
		panic("not ->")
	}

	return Arrow, "->"
}

// Scan reads the buffer and returns the peak token and literal.
func (l *lexer) Scan() (token, string) {
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
