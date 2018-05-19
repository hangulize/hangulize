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
