package hangulize

import (
	"embed"
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
func LoadSpec(lang string) (*Spec, bool) {
	var spec *Spec

	spec, ok := specs[lang]
	if ok {
		// already loaded
		return spec, true
	}

	filename := "specs/" + lang + ext
	hsl, err := f.ReadFile(filename)

	if errors.Is(err, fs.ErrNotExist) {
		// not found
		return nil, false
	}

	spec, err = ParseSpec(strings.NewReader(string(hsl)))

	// Bundled spec must not have any error.
	if err != nil {
		panic(errors.Wrapf(err, `bundled spec "%s" has error`, lang))
	}

	// Cache it.
	specs[lang] = spec
	return spec, true
}

// UnloadSpec flushes a cached spec to get free memory.
func UnloadSpec(lang string) {
	delete(specs, lang)
}
