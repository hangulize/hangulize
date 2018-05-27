package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeHangul(t *testing.T) {
	assert.Equal(t, "한글", ComposeHangul("ㅎㅏ-ㄴㄱㅡ-ㄹ"))
	assert.Equal(t, "낑깡", ComposeHangul("ㄲㅣ-ㅇㄲㅏ-ㅇ"))
}
