package main

import (
	"bytes"
	"fmt"
	"testing"
)

func BenchmarkRun(b *testing.B) {
	b.Run("naive", func(b *testing.B) {
		var buf bytes.Buffer
		for i := 0; i < b.N; i++ {
			buf.Reset()
			naive(&buf)
		}
	})

	for n := 1; n <= 32; n++ {
		b.Run(fmt.Sprintf("N=%d", n), func(b *testing.B) {
			var buf bytes.Buffer
			for i := 0; i < b.N; i++ {
				buf.Reset()
				run(&buf, n)
			}
		})
	}
}
