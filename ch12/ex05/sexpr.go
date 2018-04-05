package sexpr

import (
	"bytes"
	"fmt"
	"reflect"
)

func encode(buf *bytes.Buffer, v reflect.Value, indent int) error {
	switch v.Kind() {
	case reflect.Invalid:
		buf.WriteString("undefined")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		fmt.Fprintf(buf, "%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		fmt.Fprintf(buf, "%d", v.Uint())
	case reflect.String:
		fmt.Fprintf(buf, "%q", v.String())
	case reflect.Ptr:
		if !v.IsNil() {
			return encode(buf, v.Elem(), indent)
		}
		buf.WriteString("null")

	case reflect.Array, reflect.Slice:
		buf.WriteByte('[')
		indent++
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				buf.WriteString(",\n")
				buf.Write(bytes.Repeat([]byte{' '}, indent))
			}
			if err := encode(buf, v.Index(i), indent); err != nil {
				return err
			}
		}
		buf.WriteByte(']')
	case reflect.Struct:
		buf.WriteByte('{')
		indent++
		for i := 0; i < v.NumField(); i++ {
			if i > 0 {
				buf.WriteString(",\n")
				buf.Write(bytes.Repeat([]byte{' '}, indent))
			}
			name := v.Type().Field(i).Name
			n, _ := fmt.Fprintf(buf, "%q: ", name)
			if err := encode(buf, v.Field(i), indent+n); err != nil {
				return err
			}
		}
		buf.WriteByte('}')
	case reflect.Map:
		buf.WriteByte('{')
		indent++
		for i, key := range v.MapKeys() {
			if i > 0 {
				buf.WriteString(",\n")
				buf.Write(bytes.Repeat([]byte{' '}, indent))
			}
			if err := encode(buf, key, indent); err != nil {
				return err
			}
			buf.WriteString(": ")
			if err := encode(buf, v.MapIndex(key), indent); err != nil {
				return err
			}
		}
		buf.WriteByte('}')
	case reflect.Bool:
		if v.Bool() {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case reflect.Float32, reflect.Float64:
		fmt.Fprintf(buf, "%v", v.Float())
	case reflect.Interface:
		return encode(buf, v.Elem(), indent)
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
