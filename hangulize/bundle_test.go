package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBundledTestCases(t *testing.T) {
	for _, lang := range ListSpecs() {
		spec, ok := LoadSpec(lang)

		assert.Truef(t, ok, `failed to load "%s" spec`, lang)

		h := NewHangulizer(spec)

		for _, testCase := range spec.Test {
			loanword := testCase.Left()
			expected := testCase.Right()[0]

			hangul := h.Hangulize(loanword)

			assert.Equal(t, expected, hangul, loanword)
		}
	}
}
