package cyrillic

import (
	"fmt"
	"strings"

	"github.com/hangulize/hangulize"
)

var (
	Bulgaria  hangulize.Translit = mustNew("bg")
	Macedonia                    = mustNew("mk")
	Russia                       = mustNew("ru")
	Ukraine                      = mustNew("ua")
	// More available: tj, mn, sr, me
)

var Ts = []hangulize.Translit{Bulgaria, Macedonia, Russia, Ukraine}

func mustNew(country string) hangulize.Translit {
	dict := loadDict()

	mapping, ok := dict[country]
	if !ok {
		panic(fmt.Errorf("cyrillic: invalid country: %s", country))
	}

	oldnew := make([]string, 0, len(mapping.ToCyrillic)*2)
	for src, dst := range mapping.ToCyrillic {
		oldnew = append(oldnew, src, dst)
	}
	repl := strings.NewReplacer(oldnew...)

	return &cyrillic{country, repl}
}

type cyrillic struct {
	country string
	repl    *strings.Replacer
}

func (c *cyrillic) Scheme() string {
	return fmt.Sprintf("cyrillic[%s]", c.country)
}

func (c *cyrillic) Transliterate(word string) (string, error) {
	return c.repl.Replace(word), nil
}
