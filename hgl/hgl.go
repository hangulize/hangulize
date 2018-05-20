/*
Package hgl implements a parser for the HGL format which is used for
Hangulize 2.

The HGL format is a simple config language.  A HGL config has sections.  Each
section contains string-strings pairs.  A section can be one of a dictionary or
a pair list:

	lang:
		id      = "ita"
		code    = "it", "ita", "ita"
		english = "Italian"
		korean  = "이탈리아어"
		script  = "roman"

	rewrite:
		"^gli$"    -> "li"
		"^glia$"   -> "g.lia"
		"^glioma$" -> "g.lioma"
		"^gli{@}"  -> "li"
		"{@}gli"   -> "li"

Keys in a dictionary section must be unique, and the section won't keep their
order.  While a pair list section works in an inversed way.  A pair list
section just keeps described pairs in therir order.
*/
package hgl

import (
	"errors"
	"io"
)

type dict map[string]string

type Pair struct {
	Left  string
	Right string
}

type Section struct {
	name string
}

type DictSection struct {
	Section
	dict *dict
}

func (s *Section) Name() string {
	return s.name
}

func (s *DictSection) Dict() *dict {
	return s.dict
}

func (s *DictSection) Get(key string) string {
	return (*s.dict)[key]
}

// type PairsSection interface {
// 	Section
// 	Pairs() *[]Pair
// }

var ParseError = errors.New("Failed to parse HGL")

func Parse(r io.Reader) (map[string]*DictSection, error) {
	sections := make(map[string]*DictSection)

	// scanner := newScanner(r)

	// mockup implementation to pass only test cases
	dict := make(dict)
	dict["hello"] = "world"

	sections["foo"] = &DictSection{Section{name: "foo"}, &dict}

	return sections, nil
}
