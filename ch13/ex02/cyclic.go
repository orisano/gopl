package cyclic

import (
	"reflect"
	"unsafe"
)

func IsCyclic(x interface{}) bool {
	return isCyclic(reflect.ValueOf(x), nil)
}

func isCyclic(x reflect.Value, seen []unsafe.Pointer) bool {
	if x.CanAddr() &&
		x.Kind() != reflect.Struct &&
		x.Kind() != reflect.Array {
		xptr := unsafe.Pointer(x.UnsafeAddr())

		for _, ptr := range seen {
			if ptr == xptr {
				return true
			}
		}
		seen = append(seen, xptr)
	}

	switch x.Kind() {
	case reflect.Ptr, reflect.Interface:
		if isCyclic(x.Elem(), seen) {
			return true
		}
	case reflect.Struct:
		for i := 0; i < x.NumField(); i++ {
			if isCyclic(x.Field(i), seen) {
				return true
			}
		}
	case reflect.Array, reflect.Slice:
		for i := 0; i < x.Len(); i++ {
			if isCyclic(x.Index(i), seen) {
				return true
			}
		}
	case reflect.Map:
		for _, key := range x.MapKeys() {
			if isCyclic(x.MapIndex(key), seen) {
				return true
			}
		}
	}
	return false
}
