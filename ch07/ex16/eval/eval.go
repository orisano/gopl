package eval

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Env map[Var]float64

type Expr interface {
	Eval(env Env) float64
	Check(vars map[Var]bool) error
	String() string
}

type Var string

func (v Var) String() string {
	return string(v)
}

func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}

func (v Var) Eval(env Env) float64 {
	return env[v]
}

type literal float64

func (l literal) String() string {
	return strconv.FormatFloat(float64(l), 'g', -1, 64)
}

func (l literal) Check(vars map[Var]bool) error {
	return nil
}

func (l literal) Eval(env Env) float64 {
	return float64(l)
}

type unary struct {
	op rune
	x  Expr
}

func (u unary) String() string {
	return string(u.op) + u.x.String()
}

func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf("unexpected unary op %q", u.op)
	}
	return u.x.Check(vars)
}

func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	}
	panic(fmt.Sprintf("unsupported unary operator: %q", u.op))
}

type binary struct {
	op   rune
	x, y Expr
}

func (b binary) String() string {
	return fmt.Sprintf("%s %s %s", b.x, string(b.op), b.y)
}

func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-*/", b.op) {
		return fmt.Errorf("unexpected binary op %q", b.op)
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	}
	panic(fmt.Sprintf("unsupported binary operator: %q", b.op))
}

var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}

type call struct {
	fn   string
	args []Expr
}

func (c call) String() string {
	var b bytes.Buffer
	b.WriteString(c.fn)
	b.WriteByte('(')
	for i, arg := range c.args {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(arg.String())
	}
	b.WriteByte(')')
	return b.String()
}

func (c call) Check(vars map[Var]bool) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s has %d args, want %d", c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	}
	panic(fmt.Sprintf("unsupported function call: %q", c.fn))
}

type min struct {
	operand []Expr
}

func (m min) Eval(env Env) float64 {
	mi := m.operand[0].Eval(env)
	for _, op := range m.operand[1:] {
		if v := op.Eval(env); mi > v {
			mi = v
		}
	}
	return mi
}

func (m min) Check(vars map[Var]bool) error {
	if len(m.operand) == 0 {
		return fmt.Errorf("too few literals")
	}
	for _, op := range m.operand {
		if err := op.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

func (m min) String() string {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, op := range m.operand {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(op.String())
	}
	b.WriteByte('}')
	return b.String()
}
