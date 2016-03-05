package encoding

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func encodeQueryInnerSingle(prefix string, data string) (string, error) {
	if prefix == "" {
		return data
	} else {
		return EncodeUrl(prefix) + "=" + EncodeUrl(data)
	}
}

func encodeQueryInnerArray(prefix string, data reflect.Value) (string, error) {
	result := []string{}
	dataLen := data.Len()
	for i := 0; i != dataLen; i++ {
		singleData := data.Index(i)
		newPrefix := prefix + "[]"
		singleResult, err := encodeQueryInner(newPrefix, singleData.Interface())
		if err != nil {
			return "", err
		}
		result = append(result, singleResult)
	}
	return strings.Join(result, "&")
}

func encodeQueryInnerMap(prefix string, data reflect.Value) (string, error) {
	result := []string{}
	dataKeys := data.MapKeys()
	for _, singleDataKey := range dataKeys {
		singleDataValue := data.MapIndex(singleDataKey)
		newPrefix := prefix + "[" + fmt.Sprintf("%v", singleDataKey.Interface()) + "]"
		singleResult, err := encodeQueryInner(newPrefix, singleDataValue.Interface())
		if err != nil {
			return "", err
		}
		result = append(result, singleResult)
	}
	return strings.Join(result, "&")
}

func encodeQueryInnerStruct(prefix string, data reflect.Value) (string, error) {
	result := []string{}
	dataKeys := data.MapKeys()
}

func encodeQueryInner(prefix string, data interface{}) (string, error) {
	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	case reflect.Float32, reflect.Float64:
	case reflect.Bool:
		return encodeQueryInnerSingle(prefix, fmt.Sprintf("%v", data))
	case reflect.Slice:
		return encodeQueryInnerArray(prefix, dataValue)
	case reflect.Map:
		return encodeQueryInnerMap(prefix, dataValue)
	default:
		return "", errors.New("invalid type " + dataType)
	}
}

func EncodeQuery(data interface{}) ([]byte, error) {
	dataMap := ArrayMappingByTagOrFirstLower(data, "query")
	return encodeQueryInner("", dataMap)
}

func DecodeQuery(data []byte, value interface{}) error {

}
