package ex19

import "testing"

func TestF(t *testing.T) {
	expected := "hello"
	if got := F(); got != expected {
		t.Errorf("unexpected result. expected: %v, but got: %v", expected, got)
	}
}
