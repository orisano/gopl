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

func BenchmarkPopCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := rand.Uint64()
		PopCount(x)
	}
}

func BenchmarkPopCountBitMagic(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := rand.Uint64()
		PopCountBitMagic(x)
	}
}

func BenchmarkPopCountParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := rand.Uint64()
		PopCountParallel(x)
	}
}

func BenchmarkPopCount16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := rand.Uint64()
		PopCount16(x)
	}
}

func BenchmarkPopCount15(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := rand.Uint64()
		PopCount15(x)
	}
}
