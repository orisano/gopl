package sexpr_test

import (
	"testing"

	"github.com/orisano/gopl/ch12/ex03"
)

func TestMarshal(t *testing.T) {
	var x interface{} = []int{1, 2, 3}
	b, err := sexpr.Marshal(&x)
	if err != nil {
		t.Fatal(err)
	}
	if got, expected := string(b), `("[]int" (1 2 3))`; got != expected {
		t.Errorf("unexpected bytes. expected: %v, but got: %v", expected, got)
	}
}
