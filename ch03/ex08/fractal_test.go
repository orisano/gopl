package main

import (
	"io/ioutil"
	"testing"
)

func BenchmarkComplex64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run(ioutil.Discard, FractalFunc(complex64Fractal))
	}
}

func BenchmarkComplex128(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run(ioutil.Discard, FractalFunc(complex128Fractal))
	}
}

func BenchmarkBigFloatComplex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run(ioutil.Discard, FractalFunc(bigFloatFractal))
	}
}

func BenchmarkBigRatComplex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		run(ioutil.Discard, FractalFunc(bigRatFractal))
	}
}
