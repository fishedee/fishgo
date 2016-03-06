package language

import (
	"reflect"
)

var TypeKind struct {
	EnumStruct
	BOOL      int `enum:"1,布尔"`
	INT       int `enum:"2,有符号整数"`
	UINT      int `enum:"3,无符号整数"`
	FLOAT     int `enum:"4,浮点数"`
	PTR       int `enum:"5,指针"`
	STRING    int `enum:"6,字符串"`
	ARRAY     int `enum:"7,数组"`
	MAP       int `enum:"8,映射"`
	STRUCT    int `enum:"9,结构体"`
	INTERFACE int `enum:"10,接口"`
	FUNC      int `enum:"11,函数"`
	CHAN      int `enum:"12,通道"`
	OTHER     int `enum:"13,其他"`
}

func init() {
	InitEnumStruct(&TypeKind)
}

func GetTypeKind(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Bool:
		return TypeKind.BOOL
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return TypeKind.INT
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return TypeKind.UINT
	case reflect.Float32, reflect.Float64:
		return TypeKind.FLOAT
	case reflect.Ptr:
		return TypeKind.PTR
	case reflect.String:
		return TypeKind.STRING
	case reflect.Array, reflect.Slice:
		return TypeKind.ARRAY
	case reflect.Map:
		return TypeKind.MAP
	case reflect.Struct:
		return TypeKind.STRUCT
	case reflect.Interface:
		return TypeKind.INTERFACE
	case reflect.Func:
		return TypeKind.FUNC
	case reflect.Chan:
		return TypeKind.CHAN
	default:
		return TypeKind.OTHER
	}
}

func IsEmptyValue(v reflect.Value) bool {
	switch GetTypeKind(v.Type()) {
	case TypeKind.ARRAY, TypeKind.MAP, TypeKind.STRING:
		return v.Len() == 0
	case TypeKind.BOOL:
		return !v.Bool()
	case TypeKind.INT:
		return v.Int() == 0
	case TypeKind.UINT:
		return v.Uint() == 0
	case TypeKind.FLOAT:
		return v.Float() == 0
	case TypeKind.INTERFACE, TypeKind.PTR:
		return v.IsNil()
	}
	return false
}
