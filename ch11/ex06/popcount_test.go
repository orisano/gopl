package popcount

import (
	"math/rand"
	"testing"
)

func TestPopCount(t *testing.T) {
	for i := uint64(0); i < 1000000; i++ {
		pc := PopCount(i)
		if got := PopCountBitMagic(i); got != pc {
			t.Errorf("unexpected BitMagic result. expected: %v, but got: %v", pc, got)
		}
		if got := PopCountParallel(i); got != pc {
			t.Errorf("unexpected Parallel result. expected: %v, but got: %v", pc, got)
		}
		if got := PopCount16(i); got != pc {
			t.Errorf("unexpected 16 result. expected: %v, but got: %v", pc, got)
		}
		if got := PopCount15(i); got != pc {
			t.Errorf("unexpected 15 result. expected: %v, but got: %v", pc, got)
		}
	}
}

func benchmarkPopCount(b *testing.B, f func(uint64) int) {
	b.Helper()

	rand.Seed(0)
	s := 0
	for i := 0; i < b.N; i++ {
		x := rand.Uint64()
		s += f(x)
	}
}

func BenchmarkPopCount(b *testing.B) {
	benchmarkPopCount(b, PopCount)
}

func BenchmarkPopCountNaive(b *testing.B) {
	benchmarkPopCount(b, PopCountNaive)
}

func BenchmarkPopCountBitMagic(b *testing.B) {
	benchmarkPopCount(b, PopCountBitMagic)
}

func BenchmarkPopCountParallel(b *testing.B) {
	benchmarkPopCount(b, PopCountParallel)
}

func BenchmarkPopCount16(b *testing.B) {
	benchmarkPopCount(b, PopCount16)
}

func BenchmarkPopCount15(b *testing.B) {
	benchmarkPopCount(b, PopCount15)
}

func BenchmarkPopCountCPU(b *testing.B) {
	benchmarkPopCount(b, PopCountCPU)
}
