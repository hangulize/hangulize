package hangulize

import (
	"strings"
	"testing"
)

func BenchmarkCappuccino(b *testing.B) {
	spec, _ := LoadSpec("ita")
	h := NewHangulizer(spec)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h.Hangulize("Cappuccino")
	}
}

func BenchmarkCappuccinoTrace(b *testing.B) {
	spec, _ := LoadSpec("ita")
	h := NewHangulizer(spec)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h.HangulizeTrace("Cappuccino")
	}
}

func BenchmarkJulianaLouiseEmmaMarieWilhelmina(b *testing.B) {
	spec, _ := LoadSpec("nld")
	h := NewHangulizer(spec)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		h.Hangulize("Juliana Louise Emma Marie Wilhelmina")
	}
}

func BenchmarkVeryLongWord(b *testing.B) {
	spec, _ := LoadSpec("deu")
	h := NewHangulizer(spec)

	hunk := "Donaudampfschifffahrtselektrizitätenhauptbetriebswerkbauunterbeamtengesellschaft"

	genFunc := func(n int) func(*testing.B) {
		w := strings.Repeat(hunk, n)

		return func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				h.Hangulize(w)
			}
		}
	}

	b.Run("1", genFunc(1))
	b.Run("10", genFunc(10))
	b.Run("100", genFunc(100))
	b.Run("1000", genFunc(1000))
	b.Run("10000", genFunc(10000))
}

func BenchmarkVeryLongNegativeLookbehind(b *testing.B) {
	spec, _ := LoadSpec("nld")
	h := NewHangulizer(spec)

	// This hunk triggers the "{~@}rj" pattern.
	hunk := "rj"

	genFunc := func(n int) func(*testing.B) {
		w := strings.Repeat(hunk, n)

		return func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				h.Hangulize(w)
			}
		}
	}

	b.Run("1", genFunc(1))
	b.Run("10", genFunc(10))
	b.Run("100", genFunc(100))
	b.Run("1000", genFunc(1000))
	b.Run("10000", genFunc(10000))
	b.Run("100000", genFunc(100000))
}
