package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/orisano/gopl/ch07/ex16/eval"
)

func parseAndCheck(s string, env eval.Env) (eval.Expr, error) {
	if s == "" {
		return nil, fmt.Errorf("empty expression")
	}
	expr, err := eval.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %q: %v", s, err)
	}
	vars := make(map[eval.Var]bool)
	if err := expr.Check(vars); err != nil {
		return nil, fmt.Errorf("semantic error: %v", err)
	}
	for v := range vars {
		if _, ok := env[v]; !ok {
			return nil, fmt.Errorf("%s is undefined", v)
		}
	}
	return expr, nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		env := eval.Env{}
		for k, vs := range req.Form {
			if k == "expr" {
				continue
			}
			v, err := strconv.ParseFloat(vs[0], 64)
			if err != nil {
				fmt.Fprintf(w, "invalid variable %s=%q\n", k, vs[0])
				continue
			}
			env[eval.Var(k)] = v
		}

		for _, exprStr := range req.Form["expr"] {
			fmt.Fprint(w, exprStr, " => ")
			expr, err := parseAndCheck(exprStr, env)
			if err != nil {
				fmt.Fprintln(w, err)
				continue
			}
			fmt.Fprintln(w, expr.Eval(env))
		}
	})
	log.Fatal(http.ListenAndServe(":8000", nil))
}
