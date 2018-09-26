package hangulize

import (
	"github.com/hangulize/hangulize/internal/scripts"
)

// script represents a writing system.
type script interface {
	Is(rune) bool
	Normalize(rune) rune
	TransliteratePunct(rune) string
}

// getScript chooses a script by the script name.
func getScript(name string) (script, bool) {
	script, ok := scriptRegistry[name]
	return script, ok
}

// scriptRegistry is the registry of Scripts by their name.
var scriptRegistry map[string]script

func init() {
	scriptRegistry = map[string]script{
		"":         scripts.Latin{},
		"cyrillic": scripts.Cyrillic{},
		"georgian": scripts.Georgian{},
		"greek":    scripts.Greek{},
		"kana":     scripts.Kana{},
		"latin":    scripts.Latin{},
		"pinyin":   scripts.Pinyin{},
	}
}
