package hgl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tok token
var lit string

func _newLexer(src string) *lexer {
	return newLexer(strings.NewReader(strings.TrimSpace(src)))
}

// scan calls Lexer.Scan() but ignores Space and Newline tokens
func scan(l *lexer) (token, string) {
	for {
		tok, lit := l.Scan()
		if tok != Space && tok != Newline {
			return tok, lit
		}
	}
}

func TestSingle(t *testing.T) {
	s := _newLexer(`
	single = foo_bar_123
	`)

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "single", lit)

	tok, lit = scan(s)
	assert.Equal(t, Equal, tok)

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "foo_bar_123", lit)
}

func TestList(t *testing.T) {
	s := _newLexer(`
	list = one, "2", "셋"
	`)

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "list", lit)

	tok, lit = scan(s)
	assert.Equal(t, Equal, tok)

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "one", lit)

	tok, lit = scan(s)
	assert.Equal(t, Comma, tok)

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "2", lit)

	tok, lit = scan(s)
	assert.Equal(t, Comma, tok)

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "셋", lit)
}

func TestEscapedQuote(t *testing.T) {
	s := _newLexer(`
	escaped = "\""
	`)

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "escaped", lit)

	tok, lit = scan(s)
	assert.Equal(t, Equal, tok)

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, `"`, lit)
}

func TestCommentSingleLine(t *testing.T) {
	s := _newLexer(`
	# Hello, world!
	`)
	tok, lit = scan(s)
	assert.Equal(t, Comment, tok)
	assert.Equal(t, "Hello, world!", lit)
}

func TestCommentMultipleLines(t *testing.T) {
	s := _newLexer(`
	# Hello,
	# world!
	`)
	tok, lit = scan(s)
	assert.Equal(t, Comment, tok)
	assert.Equal(t, "Hello, world!", lit)
}

func TestCommentParagraphs(t *testing.T) {
	s := _newLexer(`
	# Hello,
	# world!
	#
	# It's the second paragraph.

	# foo
	# bar
	#
	# baz
	#
	#
	# qux
	`)

	tok, lit = scan(s)
	assert.Equal(t, Comment, tok)
	assert.Equal(t, "Hello, world!\n\nIt's the second paragraph.", lit)

	tok, lit = scan(s)
	assert.Equal(t, Comment, tok)
	assert.Equal(t, "foo bar\n\nbaz\n\nqux", lit)
}

func TestSharpInString(t *testing.T) {
	s := _newLexer(`
	"#sharp#" # comment
	`)

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "#sharp#", lit)

	tok, lit = scan(s)
	assert.Equal(t, Comment, tok)
	assert.Equal(t, "comment", lit)
}

func TestSimpleComplete(t *testing.T) {
	s := _newLexer(`
	section1:
		hello = world
		"foo" = "bar baz"
	
	section2:
		a -> b
		b -> a
		a -> c
	`)

	// seciton1:

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "section1", lit)

	tok, lit = scan(s)
	assert.Equal(t, Colon, tok)

	// hello = world

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "hello", lit)

	tok, lit = scan(s)
	assert.Equal(t, Equal, tok)

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "world", lit)

	// "foo" = "bar baz"

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "foo", lit)

	tok, lit = scan(s)
	assert.Equal(t, Equal, tok)

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "bar baz", lit)

	// seciton2:

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "section2", lit)

	tok, lit = scan(s)
	assert.Equal(t, Colon, tok)

	// a -> b

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "a", lit)

	tok, lit = scan(s)
	assert.Equal(t, Arrow, tok)

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "b", lit)

	// b -> a

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "b", lit)

	tok, lit = scan(s)
	assert.Equal(t, Arrow, tok)

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "a", lit)

	// a -> c

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "a", lit)

	tok, lit = scan(s)
	assert.Equal(t, Arrow, tok)

	tok, lit = scan(s)
	assert.Equal(t, String, tok)
	assert.Equal(t, "c", lit)
}
