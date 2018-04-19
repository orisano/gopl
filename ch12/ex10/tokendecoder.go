package sexpr

import (
	"strconv"
	"text/scanner"
)

type Token interface{}

type Symbol struct {
	Name string
}

type String struct {
	Value string
}

type Int struct {
	Value int
}

type StartList struct{}
type EndList struct{}

type TokenDecoder struct {
	lex *lexer
}

func (d *TokenDecoder) Token() Token {
	switch d.lex.token {
	case scanner.Ident:
		name := d.lex.text()
		d.lex.next()
		return Symbol{Name: name}
	case scanner.String:
		value, _ := strconv.Unquote(d.lex.text())
		d.lex.next()
		return String{Value: value}
	case scanner.Int:
		value, _ := strconv.ParseInt(d.lex.text(), 10, 64)
		d.lex.next()
		return Int{Value: int(value)}
	case '(':
		d.lex.next()
		return StartList{}
	case ')':
		d.lex.next()
		return EndList{}
	}
	panic("unknown token: " + d.lex.text())
}
