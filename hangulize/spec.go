package hangulize

import (
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/sublee/hangulize2/hgl"
)

// Spec represents a transactiption specification for a language.
type Spec struct {
	lang   Language
	config Config

	vars      map[string][]string
	rewrite   []hgl.Pair
	hangulize []hgl.Pair

	test []hgl.Pair
}

// Language identifies a natural language.
type Language struct {
	id      string
	code    []string
	english string
	korean  string
	script  string
}

// Config keeps some configurations for a transactiption specification.
type Config struct {
	authors []string
	stage   string
	markers []rune
}

// ParseSpec parses a Spec from an HGL source.
func ParseSpec(r io.Reader) (*Spec, error) {
	h, err := hgl.Parse(r)

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse HGL source")
	}

	return buildSpec(h)
}

func buildSpec(h hgl.HGL) (*Spec, error) {
	var sec hgl.Section
	var list *hgl.ListSection
	var dict *hgl.DictSection
	var ok bool

	// lang

	sec, ok = h["lang"]
	if !ok {
		return nil, errors.New("'lang' section required")
	}
	dict = sec.(*hgl.DictSection)
	lang := Language{
		id:      dict.One("id"),
		code:    dict.All("code"),
		english: dict.One("english"),
		korean:  dict.One("korean"),
		script:  dict.One("script"),
	}

	// config

	sec, ok = h["config"]
	if !ok {
		return nil, errors.New("'config' section required")
	}
	dict = sec.(*hgl.DictSection)

	// a marker must be 1-character.
	stringMarkers := dict.All("markers")
	markers := make([]rune, len(stringMarkers))

	for i, stringMarker := range stringMarkers {
		if len(stringMarker) != 1 {
			err := fmt.Errorf("marker %#v must be 1-character", stringMarker)
			return nil, err
		}
		markers[i] = rune(stringMarker[0])
	}

	config := Config{
		authors: dict.All("authors"),
		stage:   dict.One("stage"),
		markers: markers,
	}

	// vars (optional)

	var vars map[string][]string
	sec, ok = h["vars"]
	if ok {
		dict = sec.(*hgl.DictSection)
		vars = dict.Map()
	}

	// rewrite

	sec, ok = h["rewrite"]
	if !ok {
		return nil, errors.New("'rewrite' section required")
	}
	list = sec.(*hgl.ListSection)
	rewrite := list.Array()

	// hangulize

	sec, ok = h["hangulize"]
	if !ok {
		return nil, errors.New("'hangulize' section required")
	}
	list = sec.(*hgl.ListSection)
	hangulize := list.Array()

	// test (optional)

	var test []hgl.Pair
	sec, ok = h["test"]
	if ok {
		list = sec.(*hgl.ListSection)
		test = list.Array()
	}

	spec := Spec{
		lang,
		config,
		vars,
		rewrite,
		hangulize,
		test,
	}
	return &spec, nil
}
