package hgl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScan(t *testing.T) {
	src := strings.TrimSpace(`
	section_name:
		key1 = value
		key2 = one, "2", "셋"
		key3 = "\""
	`)

	var tok token
	var lit string

	scanner := newScanner(strings.NewReader(src))

	// call scanner.Scan() but skip Space and Newline tokens
	scan := func() (token, string) {
		for {
			tok, lit := scanner.Scan()
			if tok != Space && tok != Newline {
				return tok, lit
			}
		}
	}

	// section_name:

	tok, lit = scan()
	assert.Equal(t, String, tok)
	assert.Equal(t, "section_name", lit)

	tok, lit = scan()
	assert.Equal(t, Colon, tok)

	// key1 = value

	tok, lit = scan()
	assert.Equal(t, String, tok)
	assert.Equal(t, "key1", lit)

	tok, lit = scan()
	assert.Equal(t, Equal, tok)

	tok, lit = scan()
	assert.Equal(t, String, tok)
	assert.Equal(t, "value", lit)

	// key2 = one, "2", "셋"

	tok, lit = scan()
	assert.Equal(t, String, tok)
	assert.Equal(t, "key2", lit)

	tok, lit = scan()
	assert.Equal(t, Equal, tok)

	tok, lit = scan()
	assert.Equal(t, String, tok)
	assert.Equal(t, "one", lit)

	tok, lit = scan()
	assert.Equal(t, Comma, tok)

	tok, lit = scan()
	assert.Equal(t, String, tok)
	assert.Equal(t, "2", lit)

	tok, lit = scan()
	assert.Equal(t, Comma, tok)

	tok, lit = scan()
	assert.Equal(t, String, tok)
	assert.Equal(t, "셋", lit)

	// key3 = "\""

	tok, lit = scan()
	assert.Equal(t, String, tok)
	assert.Equal(t, "key3", lit)

	tok, lit = scan()
	assert.Equal(t, Equal, tok)

	tok, lit = scan()
	assert.Equal(t, String, tok)
	assert.Equal(t, `"`, lit)
}
