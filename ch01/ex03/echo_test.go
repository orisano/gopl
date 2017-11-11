package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func setupSlice(n int) []string {
	s := make([]string, 0, n)
	for i := 0; i < n; i++ {
		s = append(s, "aaa")
	}
	return s
}

func TestFastEcho(t *testing.T) {
	s := []string{"./fast_echo", "h", "o", "g", "e", "hoge"}
	expected := "h o g e hoge\n"

	var b bytes.Buffer
	FastEcho(&b, s)
	got := b.String()

	if got != expected {
		t.Errorf("unexpected output. expected: %v, got: %v", expected, got)
	}
}

func TestNaiveEcho(t *testing.T) {
	s := []string{"./naive_echo", "h", "o", "g", "e", "hoge"}
	expected := "h o g e hoge\n"

	var b bytes.Buffer
	NaiveEcho(&b, s)
	got := b.String()

	if got != expected {
		t.Errorf("unexpected output. expected: %v, got: %v", expected, got)
	}
}

func BenchmarkNaiveEcho(b *testing.B) {
	s := setupSlice(1000)
	for i := 0; i < b.N; i++ {
		NaiveEcho(ioutil.Discard, s)
	}
}

func BenchmarkFastEcho(b *testing.B) {
	s := setupSlice(1000)
	for i := 0; i < b.N; i++ {
		FastEcho(ioutil.Discard, s)
	}
}
