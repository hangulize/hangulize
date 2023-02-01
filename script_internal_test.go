package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestISO15924(t *testing.T) {
	for script := range scriptRegistry {
		if script == "" {
			continue
		}

		_, err := language.ParseScript(script)
		assert.NoError(t, err, "%s is not defined in ISO 15924", script)
	}
}
