package sexpr

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

func encode(buf *bytes.Buffer, v reflect.Value, indent int) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("nil")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(buf, "%d", v.Uint())
	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())
	case reflect.Ptr:
		return encode(buf, v.Elem(), indent)
	case reflect.Array, reflect.Slice:
		buf.WriteByte('(')
		indent++
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteByte('\n')
				buf.Write(bytes.Repeat([]byte{' '}, indent))
			}
			if err := encode(buf, v.Index(i), indent); err != nil {
				return err
			}
		}
		buf.WriteByte(')')
	case reflect.Struct:
		buf.WriteByte('(')
		indent++
		for i, cnt := 0, 0; i < v.NumField(); i++ {
			field := v.Field(i)
			zero := reflect.Zero(field.Type())
			if reflect.DeepEqual(field.Interface(), zero.Interface()) {
				continue
			}

			if cnt > 0 {
				buf.WriteByte('\n')
				buf.Write(bytes.Repeat([]byte{' '}, indent))
			}
			cnt++
			n, _ := fmt.Fprintf(buf, "(%s ", v.Type().Field(i).Name)
			if err := encode(buf, v.Field(i), indent+n); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	case reflect.Map:
		buf.WriteByte('(')
		indent++
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteByte('\n')
				buf.Write(bytes.Repeat([]byte{' '}, indent))
			}
			buf.WriteByte('(')
			if err := encode(buf, key, indent); err != nil {
				return err
			}
			buf.WriteByte(' ')
			if err := encode(buf, v.MapIndex(key), indent); err != nil {
				return err
			}
			buf.WriteByte(')')
		}
		buf.WriteByte(')')
	case reflect.Bool:
		if v.Bool() {
			buf.WriteByte('t')
		} else {
			buf.WriteString("nil")
		}
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%v", v.Float())
	case reflect.Complex64, reflect.Complex128:
		z := v.Complex()
		fmt.Fprintf(buf, "#C(%v %v)", real(z), imag(z))
	case reflect.Interface:
		buf.WriteByte('(')
		indent++
		n, _ := buf.WriteString(strconv.Quote(v.Elem().Type().String()))
		indent += n
		buf.WriteByte(' ')
		indent++
		encode(buf, v.Elem(), indent)
		buf.WriteByte(')')
	default:
		return fmt.Errorf("unsupported type: %s", v.Type())
	}
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := encode(&buf, reflect.ValueOf(v), 0); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
