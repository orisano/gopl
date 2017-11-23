package popcount

import (
	"math/rand"
	"testing"
)

func BenchmarkPopCount(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := rand.Uint64()
		PopCount(x)
	}
}

func BenchmarkPopCountNaive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		x := rand.Uint64()
		PopCountNaive(x)
	}
}
