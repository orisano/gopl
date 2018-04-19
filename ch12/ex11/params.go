package params

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func Unpack(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	fields := make(map[string]reflect.Value)
	v := reflect.ValueOf(ptr).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = v.Field(i)
	}

	for name, values := range req.Form {
		f := fields[name]
		if !f.IsValid() {
			continue
		}
		for _, value := range values {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))
			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}

func Pack(i interface{}) string {
	u := url.Values{}

	v := reflect.ValueOf(i).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		f := v.Field(i)
		switch f.Kind() {
		case reflect.Int, reflect.String, reflect.Bool:
			u.Add(name, format(f))
		case reflect.Array, reflect.Slice:
			for j := 0; j < f.Len(); j++ {
				u.Add(name, format(f.Index(j)))
			}
		default:
			panic("unsupported kind")
		}
	}
	return u.Encode()
}

func format(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Int:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.String:
		return v.String()
	case reflect.Bool:
		return strconv.FormatBool(v.Bool())
	}
	panic("unsupported kind")
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)
	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}
	return nil
}
