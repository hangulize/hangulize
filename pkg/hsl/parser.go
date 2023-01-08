package hsl

import (
	"errors"
	"io"
)

// parser ...
type parser struct {
	lexer *lexer

	buf struct {
		tok   token
		lit   string
		line  int
		reuse bool
	}
}

// newParser ...
func newParser(r io.Reader) *parser {
	return &parser{lexer: newLexer(r)}
}

func (p *parser) scan() (token, string, int) {
	// If unscan() performed, reuses the latest token and literal.
	if p.buf.reuse {
		p.buf.reuse = false
		return p.buf.tok, p.buf.lit, p.buf.line
	}

	// Scan the next one.
	tok, lit := p.lexer.Scan()
	line := p.lexer.Line()

	// Keep the latest token and literal to reuse.
	p.buf.tok, p.buf.lit, p.buf.line = tok, lit, line

	return tok, lit, line
}

func (p *parser) unscan() error {
	if p.buf.reuse {
		return errors.New("already unscanned once")
	}

	p.buf.reuse = true
	return nil
}

// parse ...
func (p *parser) parse() (HSL, error) {
	hsl := make(HSL)

	var (
		lastString  string
		sectionName string
		sectionLine int
	)

	var (
		tok  token
		lit  string
		line int
	)

	for {
		tok, lit, line = p.scan()

		// The common behavior for useless tokens.
		if tok == Illegal {
			return nil, illegalError(lit)
		} else if tok == EOF {
			break
		} else if tok == Comment {
			continue
		}

		// Remember the last string. It will be a section name or a key.
		if tok == String {
			lastString = lit
			continue
		}

		// If a colon found, the last string is a section name.
		if tok == Colon {
			sectionName = lastString
			sectionLine = line
			continue
		}

		if tok == Equal || tok == Arrow {
			if sectionName == "" {
				return nil, errors.New("pair found not in section")
			}

			values, err := p.parseValues()
			if err != nil {
				return nil, err
			}

			var section Section
			var ok bool

			// If an equals sign found, the last string is a key in a dict
			// section.
			if tok == Equal {
				section, ok = hsl[sectionName]

				if !ok {
					section = newDictSection(sectionLine)
					hsl[sectionName] = section
				}
			}

			// If an arrow sign found, the last string is a left value in a
			// pairs section.
			if tok == Arrow {
				section, ok = hsl[sectionName]

				if !ok {
					section = newListSection(sectionLine)
					hsl[sectionName] = section
				}
			}

			section.addPair(lastString, values, line)
			continue
		}
	}

	return hsl, nil
}

func (p *parser) parseValues() ([]string, error) {
	values := make([]string, 0)

	for {
		tok, lit, _ := p.scan()

		// The common behavior for useless tokens.
		if tok == Illegal {
			return nil, illegalError(lit)
		} else if tok == EOF {
			break
		} else if tok == Comment {
			continue
		}

		// Collect strings in values.
		if tok == String {
			values = append(values, lit)
		}

		// There is a more value.
		if tok == Comma {
			continue
		}

		// Values cannot be written over multiple lines.
		if tok == Newline {
			break
		}
	}

	return values, nil
}
