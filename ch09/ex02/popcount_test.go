package popcount

import (
	"math/rand"
	"testing"
)

func BenchmarkPopCount(b *testing.B) {
	rand.Seed(72)
	for i := 0; i < b.N; i++ {
		x := rand.Uint64()
		PopCount(x)
	}
}

func BenchmarkPopCountLoop(b *testing.B) {
	rand.Seed(72)
	for i := 0; i < b.N; i++ {
		x := rand.Uint64()
		PopCountLoop(x)
	}
}
