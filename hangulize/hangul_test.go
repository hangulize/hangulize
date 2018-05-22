package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompleteHangul(t *testing.T) {
	assert.Equal(t, "한글", CompleteHangul("ㅎㅏ-ㄴㄱㅡ-ㄹ"))
}
