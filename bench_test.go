package hangulize_test

import (
	"strings"
	"testing"

	"github.com/hangulize/hangulize"
)

func BenchmarkCappuccino(b *testing.B) {
	spec, _ := hangulize.LoadSpec("ita")
	h := hangulize.New(spec)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = h.Hangulize("Cappuccino")
	}
}

func BenchmarkCappuccinoTrace(b *testing.B) {
	spec, _ := hangulize.LoadSpec("ita")
	h := hangulize.New(spec)

	h.Trace(func(hangulize.Trace) {})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = h.Hangulize("Cappuccino")
	}
}

func BenchmarkJulianaLouiseEmmaMarieWilhelmina(b *testing.B) {
	spec, _ := hangulize.LoadSpec("nld")
	h := hangulize.New(spec)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = h.Hangulize("Juliana Louise Emma Marie Wilhelmina")
	}
}

func BenchmarkVeryLongWord(b *testing.B) {
	spec, _ := hangulize.LoadSpec("deu")
	h := hangulize.New(spec)

	hunk := "Donaudampfschifffahrtselektrizitätenhauptbetriebswerkbauunterbeamtengesellschaft"

	genFunc := func(n int) func(*testing.B) {
		w := strings.Repeat(hunk, n)

		return func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = h.Hangulize(w)
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
	spec, _ := hangulize.LoadSpec("nld")
	h := hangulize.New(spec)

	// This hunk triggers the "{~@}rj" pattern.
	hunk := "rj"

	genFunc := func(n int) func(*testing.B) {
		w := strings.Repeat(hunk, n)

		return func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = h.Hangulize(w)
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
