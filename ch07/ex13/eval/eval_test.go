package eval

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		expr string
		env  Env
		want string
	}{
		{"sqrt(A / pi)", Env{"A": 87616, "pi": math.Pi}, "167"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 12, "y": 1}, "1729"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
		{"5 / 9 * (F - 32)", Env{"F": -40}, "-40"},
		{"5 / 9 * (F - 32)", Env{"F": 32}, "0"},
		{"5 / 9 * (F - 32)", Env{"F": 212}, "100"},
	}
	var prevExpr string
	for _, test := range tests {
		if test.expr != prevExpr {
			fmt.Printf("\n%s\n", test.expr)
			prevExpr = test.expr
		}
		expr, err := Parse(test.expr)
		if err != nil {
			t.Error(err)
			continue
		}

		got := fmt.Sprintf("%.6g", expr.Eval(test.env))
		fmt.Printf("\t%v => %s\n", test.env, got)
		if got != test.want {
			t.Errorf("%s.Eval() in %v = %q, want %q\n", test.expr, test.env, got, test.want)
		}
	}
}

func TestString(t *testing.T) {
	tests := []string{
		"x",
		"0.2",
		"-X",
		"x+0.2",
		"x*y",
		"sqrt(x)",
		"pow(2,n)",
		"3+2*x",
		"sin(-x)*pow(1.5,-r)",
		"pow(2,sin(y))*pow(2,sin(x))/12",
		"sin(x*y/10)/10",
	}

	for _, test := range tests {
		expr, err := Parse(test)
		if err != nil {
			t.Error(err)
			continue
		}
		s := expr.String()
		got, err := Parse(s)
		if err != nil {
			t.Errorf("failed to parse %q: %v", s, err)
			continue
		}

		if !reflect.DeepEqual(got, expr) {
			t.Errorf("unexpected expr. expected: %s, but got: %s", expr, got)
		}
	}
}
