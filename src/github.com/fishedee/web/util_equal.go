package web

import (
	"fmt"
	. "github.com/fishedee/util"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

// During deepValueEqual, must keep track of checks that are
// in progress. The comparison algorithm assumes that all
// checks in progress are true when it reencounters them.
// Visited comparisons are stored in a map indexed by visit.
type visit struct {
	a1  unsafe.Pointer
	a2  unsafe.Pointer
	typ reflect.Type
}

func checkImageEqual(left string, right string) string {
	var err error
	var leftData []byte
	var rightData []byte
	if left == right {
		return ""
	}
	errorString := ""
	err = DefaultAjaxPool.Get(&Ajax{
		Url:          left,
		ResponseData: &leftData,
	})
	if err != nil && errorString == "" {
		errorString = fmt.Sprintf("Get Image %v Error %v", left, err.Error())
	}
	err = DefaultAjaxPool.Get(&Ajax{
		Url:          right,
		ResponseData: &rightData,
	})
	if err != nil && errorString == "" {
		errorString = fmt.Sprintf("Get Image %v Error %v", right, err.Error())
	}
	if reflect.DeepEqual(leftData, rightData) == false && errorString == "" {
		errorString = fmt.Sprintf("%v Image != %v Image", left, right)
	}
	return errorString
}

// Tests for deep equality using reflected types. The map argument tracks
// comparisons that have already been seen, which allows short circuiting on
// recursive types.
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
			return true
		}
		for i, n := 0, v1.NumField(); i < n; i++ {
			fieldInfo := v1.Type().Field(i)
			//xorm 特殊处理
			xormTag := fieldInfo.Tag.Get("xorm")
			xormTagInfo := strings.Split(xormTag, " ")
			isImage := false
			for _, singleXormTagInfo := range xormTagInfo {
				if singleXormTagInfo == "image" {
					isImage = true
					break
				}
			}
			if isImage {
				equalResult := checkImageEqual(v1.Field(i).String(), v2.Field(i).String())
				if equalResult != "" {
					*errorString = fmt.Sprintf("%v: %v", equalDesc, equalResult)
					return false
				} else {
					continue
				}
			}
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
			if !val1.IsValid() || !val2.IsValid() || !deepValueEqual(v1.MapIndex(k), v2.MapIndex(k), visited, depth+1, errorString, fmt.Sprintf("%v=>%v", equalDesc, k)) {
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

// DeepEqual reports whether x and y are ``deeply equal,'' defined as follows.
// Two values of identical type are deeply equal if one of the following cases applies.
// Values of distinct types are never deeply equal.
//
// Array values are deeply equal when their corresponding elements are deeply equal.
//
// Struct values are deeply equal if their corresponding fields,
// both exported and unexported, are deeply equal.
//
// Func values are deeply equal if both are nil; otherwise they are not deeply equal.
//
// Interface values are deeply equal if they hold deeply equal concrete values.
//
// Map values are deeply equal if they are the same map object
// or if they have the same length and their corresponding keys
// (matched using Go equality) map to deeply equal values.
//
// Pointer values are deeply equal if they are equal using Go's == operator
// or if they point to deeply equal values.
//
// Slice values are deeply equal when all of the following are true:
// they are both nil or both non-nil, they have the same length,
// and either they point to the same initial entry of the same underlying array
// (that is, &x[0] == &y[0]) or their corresponding elements (up to length) are deeply equal.
// Note that a non-nil empty slice and a nil slice (for example, []byte{} and []byte(nil))
// are not deeply equal.
//
// Other values - numbers, bools, strings, and channels - are deeply equal
// if they are equal using Go's == operator.
//
// In general DeepEqual is a recursive relaxation of Go's == operator.
// However, this idea is impossible to implement without some inconsistency.
// Specifically, it is possible for a value to be unequal to itself,
// either because it is of func type (uncomparable in general)
// or because it is a floating-point NaN value (not equal to itself in floating-point comparison),
// or because it is an array, struct, or interface containing
// such a value.
// On the other hand, pointer values are always equal to themselves,
// even if they point at or contain such problematic values,
// because they compare equal using Go's == operator, and that
// is a sufficient condition to be deeply equal, regardless of content.
// DeepEqual has been defined so that the same short-cut applies
// to slices and maps: if x and y are the same slice or the same map,
// they are deeply equal regardless of content.
func DeepEqual(x, y interface{}) (string, bool) {
	if x == nil || y == nil {
		return "nil != nonil", x == y
	}
	v1 := reflect.ValueOf(x)
	v2 := reflect.ValueOf(y)
	if v1.Type() != v2.Type() {
		return fmt.Sprintf("%v type != %v type", v1.String(), v2.String()), false
	}
	var errorString string
	equalDesc := ""
	result := deepValueEqual(v1, v2, make(map[visit]bool), 0, &errorString, equalDesc)
	return errorString, result
}
