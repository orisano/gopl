package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/orisano/gopl/ch07/ex15/eval"
)

func parseAndCheck(s string) (eval.Expr, map[eval.Var]bool, error) {
	if s == "" {
		return nil, nil, fmt.Errorf("empty expression")
	}
	expr, err := eval.Parse(s)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse %q: %v", s, err)
	}
	vars := make(map[eval.Var]bool)
	if err := expr.Check(vars); err != nil {
		return nil, nil, fmt.Errorf("semantic error: %v", err)
	}
	return expr, vars, nil
}

func main() {
	in := bufio.NewScanner(os.Stdin)

	var (
		expr eval.Expr
		vars map[eval.Var]bool
	)
	for {
		fmt.Print(">> ")
		if !in.Scan() {
			log.Fatal("unexpected EOF")
		}

		var err error
		expr, vars, err = parseAndCheck(in.Text())
		if err == nil {
			break
		}
		fmt.Fprintln(os.Stderr, err)
	}

	env := eval.Env{}
	for v := range vars {
		for {
			fmt.Print(v, " = ")
			if !in.Scan() {
				log.Fatal("unexpected EOF")
			}

			value, err := strconv.ParseFloat(in.Text(), 64)
			if err == nil {
				env[v] = value
				break
			}
			fmt.Fprintln(os.Stderr, err)
		}
	}

	fmt.Println("=>", expr.Eval(env))
}
