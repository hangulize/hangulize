package jamo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComposeHangul(t *testing.T) {
	assert.Equal(t, "한글", ComposeHangul("ㅎㅏ-ㄴㄱㅡ-ㄹ"))
	assert.Equal(t, "낑깡", ComposeHangul("ㄲㅣ-ㅇㄲㅏ-ㅇ"))
}

func TestComposeHangulOnComposed(t *testing.T) {
	assert.Equal(t, "한글", ComposeHangul("한글"))
	assert.Equal(t, "한글라이즈", ComposeHangul("하-ㄴ글ㄹㅏ이ㅈ"))
}

func TestComposeHangulNonHangul(t *testing.T) {
	assert.Equal(t, "Hello, world", ComposeHangul("Hello, world"))
	assert.Equal(t, "안녕, world", ComposeHangul("ㅇㅏ-ㄴㄴㅕ-ㅇ, world"))
}

// -----------------------------------------------------------------------------
// Benchmarks

func BenchmarkComposeHangul(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ComposeHangul("ㅎㅏ-ㄴㄱㅡ-ㄹ")
	}
}

// -----------------------------------------------------------------------------
// Examples

func ExampleComposeHangul_perfect() {
	fmt.Println(ComposeHangul("ㅎㅏ-ㄴㄱㅡ-ㄹㄹㅏㅇㅣㅈㅡ"))
	// Output: 한글라이즈
}

func ExampleComposeHangul_interpolation() {
	fmt.Println(ComposeHangul("ㅗㅈ"))
	// Output: 오즈
}
