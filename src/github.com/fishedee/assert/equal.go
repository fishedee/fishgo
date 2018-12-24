package assert

import (
	"fmt"
	"reflect"
	"time"
	"unsafe"
)

type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

func deepValueEqual(v1, v2 reflect.Value, visited map[visit]bool, depth int, errorString *string, equalDesc string) bool {
	if !v1.IsValid() || !v2.IsValid() {
		*errorString = fmt.Sprintf("%v: valid != novalid", equalDesc)
		return v1.IsValid() == v2.IsValid()
	}
	if v1.Type() != v2.Type() {
		*errorString = fmt.Sprintf("%v: %v type != %v type", equalDesc, v1.Type(), v2.Type())
		return false
	}

	// if depth > 10 { panic("deepValueEqual") }	// for debugging
	hard := func(k reflect.Kind) bool {
		switch k {
		case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
			return true
		}
		return false
	}

	if v1.CanAddr() && v2.CanAddr() && hard(v1.Kind()) {
		addr1 := unsafe.Pointer(v1.UnsafeAddr())
		addr2 := unsafe.Pointer(v2.UnsafeAddr())
		if uintptr(addr1) > uintptr(addr2) {
			// Canonicalize order to reduce number of entries in visited.
			// Assumes non-moving garbage collector.
			addr1, addr2 = addr2, addr1
		}

		// Short circuit if references are already seen.
		typ := v1.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return true
		}

		// Remember for later.
		visited[v] = true
	}

	switch v1.Kind() {
	case reflect.Array:
		for i := 0; i < v1.Len(); i++ {
			if !deepValueEqual(v1.Index(i), v2.Index(i), visited, depth+1, errorString, fmt.Sprintf("%v=>%v", equalDesc, i)) {
				return false
			}
		}
		return true
	case reflect.Slice:
		if v1.Len() != v2.Len() {
			*errorString = fmt.Sprintf("%v: len(slice)[%v] != len(slice)[%v]", equalDesc, v1.Len(), v2.Len())
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for i := 0; i < v1.Len(); i++ {
			if !deepValueEqual(v1.Index(i), v2.Index(i), visited, depth+1, errorString, fmt.Sprintf("%v=>%v", equalDesc, i)) {
				return false
			}
		}
		return true
	case reflect.Interface:
		if v1.IsNil() || v2.IsNil() {
			*errorString = fmt.Sprintf("%v: %#v != %#v", equalDesc, v1, v2)
			return v1.IsNil() == v2.IsNil()
		}
		return deepValueEqual(v1.Elem(), v2.Elem(), visited, depth+1, errorString, equalDesc)
	case reflect.Ptr:
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		return deepValueEqual(v1.Elem(), v2.Elem(), visited, depth+1, errorString, equalDesc)
	case reflect.Struct:
		_, isTime := v1.Interface().(time.Time)
		if isTime {
			time1 := v1.Interface().(time.Time)
			time2 := v2.Interface().(time.Time)
			if !time1.Equal(time2) {
				*errorString = fmt.Sprintf("%v: %#v != %#v", equalDesc, time1.Format("2006-01-02 15:04:05"), time2.Format("2006-01-02 15:04:05"))
				return false
			}
			return true
		}
		for i, n := 0, v1.NumField(); i < n; i++ {
			fieldInfo := v1.Type().Field(i)
			if !deepValueEqual(v1.Field(i), v2.Field(i), visited, depth+1, errorString, fmt.Sprintf("%v=>%v", equalDesc, fieldInfo.Name)) {
				return false
			}
		}
		return true
	case reflect.Map:
		if v1.Len() != v2.Len() {
			*errorString = fmt.Sprintf("%v: len(map)[%v] != len(map)[%v]", equalDesc, v1.Len(), v2.Len())
			return false
		}
		if v1.Pointer() == v2.Pointer() {
			return true
		}
		for _, k := range v1.MapKeys() {
			val1 := v1.MapIndex(k)
			val2 := v2.MapIndex(k)
			if !val1.IsValid() || !val2.IsValid() {
				*errorString = fmt.Sprintf("%v=>%v: exist != noexist", equalDesc, k)
				return false
			}
			if !deepValueEqual(v1.MapIndex(k), v2.MapIndex(k), visited, depth+1, errorString, fmt.Sprintf("%v=>%v", equalDesc, k)) {
				return false
			}
		}
		return true
	case reflect.Func:
		if v1.IsNil() && v2.IsNil() {
			return true
		}
		// Can't do better than this:
		*errorString = fmt.Sprintf("%v: function can't not equal", equalDesc)
		return false
	default:
		// Normal equality suffices
		*errorString = fmt.Sprintf("%v: %#v != %#v", equalDesc, v1.Interface(), v2.Interface())
		return reflect.DeepEqual(v1.Interface(), v2.Interface())
	}
}

func DeepEqual(x, y interface{}) (string, bool) {
	if x == nil || y == nil {
		if x == y {
			return "", true
		} else {
			return "nil != nonil", false
		}
	}
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	if v1.Type() != v2.Type() {
		return fmt.Sprintf("%v type != %v type", v1.String(), v2.String()), false
	}
	var errorString string
	equalDesc := ""
	result := deepValueEqual(v1, v2, make(map[visit]bool), 0, &errorString, equalDesc)
	if result == true {
		errorString = ""
	}
	return errorString, result
}
