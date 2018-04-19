package sexpr

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"
)

type lexer struct {
	scan  scanner.Scanner
	token rune
}

func (lex *lexer) next()        { lex.token = lex.scan.Scan() }
func (lex *lexer) text() string { return lex.scan.TokenText() }

func (lex *lexer) consume(want rune) {
	if lex.token != want {
		panic(fmt.Sprintf("got %q, want %q", lex.text(), want))
	}
	lex.next()
}

func readFloat(lex *lexer) float64 {
	if lex.token != scanner.Float {
		panic("invalid complex:" + lex.text())
	}
	f, _ := strconv.ParseFloat(lex.text(), 64)
	return f
}

func read(lex *lexer, v reflect.Value) {
	switch lex.token {
	case scanner.Ident:
		if lex.text() == "nil" {
			v.Set(reflect.Zero(v.Type()))
			lex.next()
			return
		}
		if lex.text() == "t" {
			v.SetBool(true)
			lex.next()
			return
		}
	case scanner.String:
		s, _ := strconv.Unquote(lex.text())
		v.SetString(s)
		lex.next()
		return
	case scanner.Int:
		i, _ := strconv.ParseInt(lex.text(), 10, 64)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v.SetInt(i)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v.SetUint(uint64(i))
		default:
			panic("invalid type. cannot set int")
		}
		lex.next()
		return

	case scanner.Float:
		v.SetFloat(readFloat(lex))
		lex.next()
		return

	case '-':
		lex.consume('-')
		switch lex.token {
		case scanner.Int:
			i, _ := strconv.ParseInt(lex.text(), 10, 64)
			switch v.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v.SetInt(-i)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				panic("cannot assign")
			default:
				panic("invalid type. cannot set int")
			}
		case scanner.Float:
			v.SetFloat(-readFloat(lex))
		default:
			panic("invalid type.")
		}
		lex.next()
		return

	case '#':
		lex.consume('#')
		lex.consume(scanner.Ident) // C
		lex.consume('(')
		r := readFloat(lex)
		lex.next()
		i := readFloat(lex)
		lex.next()
		lex.consume(')')
		v.SetComplex(complex(r, i))
		return

	case '(':
		lex.next()
		readList(lex, v)
		lex.next()
		return
	}
	panic(fmt.Sprintf("unexpected token %q", lex.text()))
}

func readList(lex *lexer, v reflect.Value) {
	switch v.Kind() {
	case reflect.Array:
		for i := 0; !endList(lex); i++ {
			read(lex, v.Index(i))
		}

	case reflect.Slice:
		for !endList(lex) {
			item := reflect.New(v.Type().Elem()).Elem()
			read(lex, item)
			v.Set(reflect.Append(v, item))
		}

	case reflect.Struct:
		for !endList(lex) {
			lex.consume('(')
			if lex.token != scanner.Ident {
				panic(fmt.Sprintf("got token %q, want field name", lex.text()))
			}
			name := lex.text()
			lex.next()
			read(lex, v.FieldByName(name))
			lex.consume(')')
		}

	case reflect.Map:
		v.Set(reflect.MakeMap(v.Type()))
		for !endList(lex) {
			lex.consume('(')
			key := reflect.New(v.Type().Key()).Elem()
			read(lex, key)
			value := reflect.New(v.Type().Elem()).Elem()
			read(lex, value)
			v.SetMapIndex(key, value)
			lex.consume(')')
		}

	case reflect.Interface:
		if lex.token != scanner.String {
			panic(fmt.Sprintf("got token %q, want type name", lex.text()))
		}
		t, _ := strconv.Unquote(lex.text())
		lex.next()
		value := reflect.New(getType(t)).Elem()
		read(lex, value)
		v.Set(value)

	default:
		panic(fmt.Sprintf("cannot decode list into %v", v.Type()))
	}
}

func endList(lex *lexer) bool {
	switch lex.token {
	case scanner.EOF:
		panic("end of file")
	case ')':
		return true
	}
	return false
}

func getType(t string) reflect.Type {
	switch t {
	case "int":
		return reflect.ValueOf(int(0)).Type()
	case "int8":
		return reflect.ValueOf(int8(0)).Type()
	case "int16":
		return reflect.ValueOf(int16(0)).Type()
	case "int32":
		return reflect.ValueOf(int32(0)).Type()
	case "int64":
		return reflect.ValueOf(int64(0)).Type()
	case "uint":
		return reflect.ValueOf(uint(0)).Type()
	case "uint8":
		return reflect.ValueOf(uint8(0)).Type()
	case "uint16":
		return reflect.ValueOf(uint16(0)).Type()
	case "uint32":
		return reflect.ValueOf(uint32(0)).Type()
	case "uint64":
		return reflect.ValueOf(uint64(0)).Type()
	case "float32":
		return reflect.ValueOf(float32(0)).Type()
	case "float64":
		return reflect.ValueOf(float64(0)).Type()
	case "complex64":
		return reflect.ValueOf(complex64(0)).Type()
	case "complex128":
		return reflect.ValueOf(complex128(0)).Type()
	case "byte":
		return reflect.ValueOf(byte(0)).Type()
	case "rune":
		return reflect.ValueOf(rune(0)).Type()
	case "string":
		return reflect.ValueOf("").Type()
	}

	if strings.HasPrefix(t, "[]") {
		return reflect.SliceOf(getType(t[2:]))
	}
	if strings.HasPrefix(t, "[") {
		i := strings.IndexByte(t, ']')
		if i < 0 {
			panic("invalid type: " + t)
		}
		size, err := strconv.ParseInt(t[1:i], 10, 64)
		if err != nil || size < 0 {
			panic("invalid array size: " + t[1:i])
		}
		return reflect.ArrayOf(int(size), getType(t[i+1:]))
	}
	if strings.HasPrefix(t, "map[") {
		cnt := 1
		p := 4
		for i, r := range t[4:] {
			if r == '[' {
				cnt++
			}
			if r == ']' {
				cnt--
				if cnt == 0 {
					p += i
					break
				}
			}
		}
		return reflect.MapOf(getType(t[4:p]), getType(t[p+1:]))
	}
	panic("unsupported type: " + t)
}
