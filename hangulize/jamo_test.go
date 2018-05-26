package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompleteHangul(t *testing.T) {
	assert.Equal(t, "한글", AssembleJamo("ㅎㅏ-ㄴㄱㅡ-ㄹ"))
	assert.Equal(t, "낑깡", AssembleJamo("ㄲㅣ-ㅇㄲㅏ-ㅇ"))
}
