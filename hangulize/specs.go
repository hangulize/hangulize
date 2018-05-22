package hangulize

import (
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

// cached specs
var specs map[string]*Spec

var specBox packr.Box

func init() {
	specs = make(map[string]*Spec)
	specBox = packr.NewBox("./specs")
}

// LoadSpec finds a bundled spec by the given language name.
// Once it loads a spec, it will cache the spec.
func LoadSpec(lang string) (*Spec, bool) {
	var spec *Spec

	spec, ok := specs[lang]
	if ok {
		// already loaded
		return spec, true
	}

	filename := lang + `.hgl`

	if !specBox.Has(filename) {
		// not found
		return nil, false
	}

	hgl := specBox.String(filename)
	spec, err := ParseSpec(strings.NewReader(hgl))

	// Bundled spec must not have any error.
	if err != nil {
		panic(errors.Wrapf(err, `bundled spec "%s" has error`, lang))
	}

	// Cache it.
	specs[lang] = spec
	return spec, true
}
