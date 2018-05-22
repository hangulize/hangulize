package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompleteHangul(t *testing.T) {
	assert.Equal(t, "한글", CompleteHangul("ㅎㅏ-ㄴㄱㅡ-ㄹ"))
	assert.Equal(t, "낑깡", CompleteHangul("ㄲㅣ-ㅇㄲㅏ-ㅇ"))
}
