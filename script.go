package hangulize

import (
	"github.com/hangulize/hangulize/internal/scripts"
)

// script represents a writing system.
type script interface {
	Is(rune) bool
	Normalize(rune) rune
	LocalizePunct(rune) string
}

// getScript chooses a script by the script name.
func getScript(name string) (script, bool) {
	script, ok := scriptRegistry[name]
	return script, ok
}

// scriptRegistry is the registry of Scripts by their name.
var scriptRegistry map[string]script

func init() {
	latn := scripts.Latn{}
	scriptRegistry = map[string]script{
		"":     latn,
		"Latn": latn,
		"Cyrl": scripts.Cyrl{},
		"Geor": scripts.Geor{},
		"Grek": scripts.Grek{},
		"Hrkt": scripts.Hrkt{},
	}
}
