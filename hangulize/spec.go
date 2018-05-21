package hangulize

import (
	"fmt"
	"io"
	"strings"

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

func (s *Spec) String() string {
	return fmt.Sprintf("<Spec lang=%s>", s.lang.id)
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

	// lang
	lang, err := buildLanguage(&h)
	if err != nil {
		return nil, err
	}

	// config (optional)
	config, err := buildConfig(&h)
	if err != nil {
		return nil, err
	}

	// vars (optional)
	vars, err := buildVars(&h)
	if err != nil {
		return nil, err
	}

	// rewrite
	rewrite, err := buildRewrite(&h)
	if err != nil {
		return nil, err
	}

	// hangulize
	hangulize, err := buildHangulize(&h)
	if err != nil {
		return nil, err
	}

	// test (optional)
	test := buildTest(&h)

	spec := Spec{
		*lang,
		*config,
		vars,
		rewrite,
		hangulize,
		test,
	}
	return &spec, nil
}

func buildLanguage(h *hgl.HGL) (*Language, error) {
	sec, ok := (*h)["lang"]
	if !ok {
		return nil, errors.New("'lang' section required")
	}
	dict := sec.(*hgl.DictSection)

	lang := Language{
		id:      dict.One("id"),
		code:    dict.All("code"),
		english: dict.One("english"),
		korean:  dict.One("korean"),
		script:  dict.One("script"),
	}
	return &lang, nil
}

func buildConfig(h *hgl.HGL) (*Config, error) {
	var config *Config

	sec, ok := (*h)["config"]
	if ok {
		dict := sec.(*hgl.DictSection)

		// a marker must be 1-character
		stringMarkers := dict.All("markers")
		markers := make([]rune, len(stringMarkers))

		for i, stringMarker := range stringMarkers {
			if len(stringMarker) != 1 {
				err := fmt.Errorf("marker %#v must be 1-character", stringMarker)
				return nil, err
			}
			markers[i] = rune(stringMarker[0])
		}

		config = &Config{
			authors: dict.All("authors"),
			stage:   dict.One("stage"),
			markers: markers,
		}
	} else {
		config = &Config{
			authors: make([]string, 0),
			stage:   "",
			markers: make([]rune, 0),
		}
	}

	return config, nil
}

func isVarName(varName string) bool {
	if varName == "@" {
		// Allow "@" exceptionally.  Usually it is used as <vowels> instead.
		return true
	}
	// Except "@", every variables should be wrapped in "<...>".
	return strings.HasPrefix(varName, "<") && strings.HasSuffix(varName, ">")
}

func buildVars(h *hgl.HGL) (map[string][]string, error) {
	var vars map[string][]string
	sec, ok := (*h)["vars"]
	if ok {
		dict := sec.(*hgl.DictSection)
		vars = dict.Map()

		for varName := range vars {
			if !isVarName(varName) {
				err := fmt.Errorf("%#v is not valid var name", varName)
				return nil, err
			}
		}
	}
	return vars, nil
}

func buildRewrite(h *hgl.HGL) ([]hgl.Pair, error) {
	sec, ok := (*h)["rewrite"]
	if !ok {
		return nil, errors.New("'rewrite' section required")
	}
	list := sec.(*hgl.ListSection)
	rewrite := list.Array()
	return rewrite, nil
}

func buildHangulize(h *hgl.HGL) ([]hgl.Pair, error) {
	sec, ok := (*h)["hangulize"]
	if !ok {
		return nil, errors.New("'hangulize' section required")
	}
	list := sec.(*hgl.ListSection)
	rewrite := list.Array()
	return rewrite, nil
}

func buildTest(h *hgl.HGL) []hgl.Pair {
	var test []hgl.Pair
	sec, ok := (*h)["test"]
	if ok {
		list := sec.(*hgl.ListSection)
		test = list.Array()
	}
	return test
}
