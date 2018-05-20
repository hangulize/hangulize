package hgl

import (
	"errors"
	"io"
)

// Parser ...
type Parser struct {
	lexer *Lexer
	buf   struct {
		token   Token
		literal string
		reuse   bool
	}
}

// NewParser ...
func NewParser(r io.Reader) *Parser {
	return &Parser{lexer: NewLexer(r)}
}

func (p *Parser) scan() (Token, string) {
	// If unscan() performed, reuses the latest token and literal.
	if p.buf.reuse {
		p.buf.reuse = false
		return p.buf.token, p.buf.literal
	}

	// Scan the next one.
	token, literal := p.lexer.Scan()

	// Keep the latest token and literal to reuse.
	p.buf.token, p.buf.literal = token, literal

	return token, literal
}

func (p *Parser) unscan() error {
	if p.buf.reuse {
		return errors.New("already unscanned once")
	}

	p.buf.reuse = true
	return nil
}

// Parse ...
func (p *Parser) Parse() (HGL, error) {
	hgl := make(HGL)

	var lastString string
	var sectionName string

	for {
		token, literal := p.scan()

		// The common behavior for useless tokens.
		if token == Illegal {
			return nil, IllegalError(literal)
		} else if token == EOF {
			break
		} else if token == Comment {
			continue
		}

		// Remember the last string.  It will be a section name or a key.
		if token == String {
			lastString = literal
			continue
		}

		// If a colon found, the last string is a section name.
		if token == Colon {
			sectionName = lastString
			continue
		}

		// If an equals sign found, the last string is a key in a dict section.
		if token == Equal {
			values, err := p.parseValues()
			if err != nil {
				return nil, err
			}

			var dict Dict

			section, ok := hgl[sectionName]
			if ok {
				dict = section.(Dict)
			} else {
				dict = make(Dict)
				hgl[sectionName] = dict
			}

			dict[lastString] = values
			continue
		}

		// If an arrow sign found, the last string is a left value in a pairs
		// section.
		if token == Arrow {
			values, err := p.parseValues()
			if err != nil {
				return nil, err
			}

			var pairs Pairs

			section, ok := hgl[sectionName]
			if ok {
				pairs = section.(Pairs)
			} else {
				pairs = make(Pairs, 0)
			}

			pairs = append(pairs, Pair{L: lastString, R: values})
			hgl[sectionName] = pairs
			continue
		}
	}

	return hgl, nil
}

func (p *Parser) parseValues() ([]string, error) {
	values := make([]string, 0)

	for {
		token, literal := p.scan()

		// The common behavior for useless tokens.
		if token == Illegal {
			return nil, IllegalError(literal)
		} else if token == EOF {
			break
		} else if token == Comment {
			continue
		}

		// Collect strings in values.
		if token == String {
			values = append(values, literal)
		}

		// There is a more value.
		if token == Comma {
			continue
		}

		// Values cannot be written over multiple lines.
		if token == Newline {
			break
		}
	}

	return values, nil
}
