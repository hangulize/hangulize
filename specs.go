//go:generate go get -v github.com/gobuffalo/packr/...
//go:generate packr -v

package hangulize

import (
	"sort"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

// The box for HSL files.
var hsls = packr.NewBox("./specs")

const ext = `.hsl`

// ListLangs returns the language name list of bundled specs.
// The bundled spec can be loaded by LoadSpec.
func ListLangs() []string {
	var langs []string

	for _, filename := range hsls.List() {
		if strings.HasSuffix(filename, ext) {
			langs = append(langs, strings.TrimSuffix(filename, ext))
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

	filename := lang + ext

	if !hsls.Has(filename) {
		// not found
		return nil, false
	}

	hsl := hsls.String(filename)
	spec, err := ParseSpec(strings.NewReader(hsl))

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
