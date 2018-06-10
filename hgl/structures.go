package hgl

import (
	"fmt"
)

// HGL is a decoding result of an HGL code.
type HGL map[string]Section

// Pair is a left-right tuple:
//
//  aa -> "ㅏ", "ㅐ"
//  ^^^^^^^^^^^^^^^^
//
type Pair struct {
	l string
	r []string
}

func (p *Pair) String() string {
	return fmt.Sprintf("Pair{%#v, %#v}", p.l, p.r)
}

// Left is a string. It is used for as keys in dict:
//
//  english = "Italian"
//  ^^^^^^^
//
// Or as left of pair:
//
//  aa -> "ㅏ", "ㅐ"
//  ^^
//
func (p Pair) Left() string {
	return p.l
}

// Right is a string array. It is used for as values in dict:
//
//  english = "Italian"
//	          ^^^^^^^^^
//
// Or as right of pair:
//
//  aa -> "ㅏ", "ㅐ"
//        ^^^^^^^^^^
//
func (p Pair) Right() []string {
	return p.r
}

// Section contains pairs.
type Section interface {
	Pairs() []Pair
	addPair(string, []string) error
}

// ListSection has an ordered list of pairs.
type ListSection struct {
	pairs []Pair
}

// DictSection has an unordered list of pairs.
// Each left of underlying pairs is unique.
type DictSection struct {
	dict map[string][]string
}

// newListSection creates an empty list section.
func newListSection() *ListSection {
	return &ListSection{make([]Pair, 0)}
}

// newDictSection creates an empty dict section.
func newDictSection() *DictSection {
	return &DictSection{make(map[string][]string)}
}

// common section methods

// Pairs returns underlying pairs as an array.
func (s *ListSection) Pairs() []Pair {
	return s.pairs
}

// Pairs returns dict key-values as an array of pairs.
func (s *DictSection) Pairs() []Pair {
	pairs := make([]Pair, len(s.dict))

	i := 0
	for l, r := range s.dict {
		pairs[i] = Pair{l, r}
		i++
	}

	return pairs
}

// addPair adds a pair into a list section. It never fails.
func (s *ListSection) addPair(l string, r []string) error {
	s.pairs = append(s.pairs, Pair{l, r})
	return nil
}

// addPair adds a pair into a dict section. If there's already a pair having
// same left, it will fails.
func (s *DictSection) addPair(l string, r []string) error {
	_, ok := s.dict[l]
	if ok {
		return fmt.Errorf("left of pair duplicated: %#v", l)
	}

	s.dict[l] = r
	return nil
}

// ListSection only methods

// Array returns the underying pair array of a list section.
func (s *ListSection) Array() []Pair {
	return s.pairs
}

// DictSection only methods

// Map returns the underying map of a dict section.
func (s *DictSection) Map() map[string][]string {
	return s.dict
}

// Injective returns the underying 1-to-1 map of a dict section.
// If some right (values) has multiple values, it returns an error.
func (s *DictSection) Injective() (map[string]string, error) {
	oneToOne := make(map[string]string, len(s.dict))

	for left, right := range s.dict {
		if len(right) != 1 {
			err := fmt.Errorf("right %#v has multiple values", right)
			return nil, err
		}
		oneToOne[left] = right[0]
	}

	return oneToOne, nil
}

// One assumes the given left (key) has only one right (values). Then returns
// the only right value.
func (s *DictSection) One(left string) string {
	right, ok := s.dict[left]

	if !ok {
		return ""
	}

	if len(right) == 0 {
		return ""
	}

	return right[0]
}

// All returns the right values.
func (s *DictSection) All(left string) []string {
	right, ok := s.dict[left]

	if !ok {
		return make([]string, 0)
	}

	return right
}
