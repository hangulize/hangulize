package hangulize

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

const ext = `.hsl`

//go:embed specs/*.hsl
var f embed.FS

// ListLangs returns the language name list of bundled specs.
// The bundled spec can be loaded by LoadSpec.
func ListLangs() []string {
	dir, err := f.ReadDir("specs")
	if err != nil {
		return nil
	}

	var langs []string
	for _, ent := range dir {
		if strings.HasSuffix(ent.Name(), ext) {
			langs = append(langs, strings.TrimSuffix(ent.Name(), ext))
		}
	}

	sort.Strings(langs)
	return langs
}

// Cached specs.
var specs = make(map[string]*Spec)

// LoadSpec finds a bundled spec by the given language name.
// Once it loads a spec, it will cache the spec.
func LoadSpec(lang string) (*Spec, error) {
	var spec *Spec

	spec, ok := specs[lang]
	if ok {
		// already loaded
		return spec, nil
	}

	filename := "specs/" + lang + ext
	hsl, err := f.ReadFile(filename)

	if errors.Is(err, fs.ErrNotExist) {
		// not found
		return nil, fmt.Errorf("%w: %s", ErrSpecNotFound, lang)
	}

	spec, err = ParseSpec(strings.NewReader(string(hsl)))
	if err != nil {
		// Bundled spec must not have any error.
		panic(fmt.Errorf("bundled spec '%s': %w", lang, err))
	}

	// Cache it.
	specs[lang] = spec
	return spec, nil
}

// UnloadSpec flushes a cached spec to get free memory.
func UnloadSpec(lang string) {
	delete(specs, lang)
}
