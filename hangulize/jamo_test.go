package hangulize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeHangul(t *testing.T) {
	assert.Equal(t, "한글", ComposeHangul("ㅎㅏ-ㄴㄱㅡ-ㄹ"))
	assert.Equal(t, "낑깡", ComposeHangul("ㄲㅣ-ㅇㄲㅏ-ㅇ"))
}

func BenchmarkComposeHangul(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ComposeHangul("ㅎㅏ-ㄴㄱㅡ-ㄹ")
	}
}
