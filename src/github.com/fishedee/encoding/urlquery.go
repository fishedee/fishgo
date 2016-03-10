package encoding

import (
	"errors"
	"fmt"
	"github.com/fishedee/language"
	"reflect"
	"strconv"
	"strings"
)

func encodeQueryInnerSingle(prefix string, data string) (string, error) {
	if prefix == "" {
		return data, nil
	} else {
		prefixEncode, err := EncodeUrl(prefix)
		if err != nil {
			return "", err
		}
		dataEncode, err := EncodeUrl(data)
		if err != nil {
			return "", err
		}
		return prefixEncode + "=" + dataEncode, nil
	}
}

func isSingleValue(data reflect.Value) bool {
	if data.IsValid() == false {
		return true
	}
	dataKind := language.GetTypeKind(data.Type())
	if dataKind == language.TypeKind.ARRAY ||
		dataKind == language.TypeKind.MAP {
		return false
	} else if dataKind == language.TypeKind.INTERFACE ||
		dataKind == language.TypeKind.PTR {
		return isSingleValue(data.Elem())
	} else {
		return true
	}
}

func encodeQueryInnerArray(prefix string, data reflect.Value) (string, error) {
	result := []string{}
	dataLen := data.Len()
	for i := 0; i != dataLen; i++ {
		var newPrefix string
		singleData := data.Index(i)
		isSingle := isSingleValue(singleData)
		if isSingle == false {
			newPrefix = prefix + "[" + strconv.Itoa(i) + "]"
		} else {
			newPrefix = prefix + "[]"
		}
		singleResult, err := encodeQueryInner(newPrefix, singleData)
		if err != nil {
			return "", err
		}
		result = append(result, singleResult)
	}
	return strings.Join(result, "&"), nil
}

func encodeQueryInnerMap(prefix string, data reflect.Value) (string, error) {
	result := []string{}
	dataKeys := data.MapKeys()
	for _, singleDataKey := range dataKeys {
		singleDataValue := data.MapIndex(singleDataKey)
		var newPrefix string
		if prefix != "" {
			newPrefix = prefix + "[" + fmt.Sprintf("%v", singleDataKey.Interface()) + "]"
		} else {
			newPrefix = fmt.Sprintf("%v", singleDataKey.Interface())
		}
		singleResult, err := encodeQueryInner(newPrefix, singleDataValue)
		if err != nil {
			return "", err
		}
		result = append(result, singleResult)
	}
	return strings.Join(result, "&"), nil
}

func encodeQueryInner(prefix string, data reflect.Value) (string, error) {
	if data.IsValid() == false {
		return "", nil
	}
	dataType := data.Type()
	dataKind := language.GetTypeKind(dataType)
	if dataKind == language.TypeKind.BOOL ||
		dataKind == language.TypeKind.INT ||
		dataKind == language.TypeKind.UINT ||
		dataKind == language.TypeKind.FLOAT ||
		dataKind == language.TypeKind.STRING {
		return encodeQueryInnerSingle(prefix, fmt.Sprintf("%v", data))
	} else if dataKind == language.TypeKind.ARRAY {
		return encodeQueryInnerArray(prefix, data)
	} else if dataKind == language.TypeKind.MAP {
		return encodeQueryInnerMap(prefix, data)
	} else if dataKind == language.TypeKind.PTR ||
		dataKind == language.TypeKind.INTERFACE {
		return encodeQueryInner(prefix, data.Elem())
	} else {
		return "", errors.New("invalid type " + dataType.String())
	}
}

func EncodeUrlQuery(data interface{}) ([]byte, error) {
	dataMap := language.ArrayToMap(data, "url")
	result, err := encodeQueryInner("", reflect.ValueOf(dataMap))
	if err != nil {
		return nil, err
	}
	return []byte(result), nil
}

func decodeQueryKey(data string) (string, string) {
	index := strings.Index(data, "[")
	if index != -1 {
		return data[0:index], data[index:]
	} else {
		return data, ""
	}
}

func decodeQueryValue(curData interface{}, prefix string, value string) (interface{}, error) {
	if prefix == "" {
		return value, nil
	}
	if prefix[0] != '[' {
		return nil, errors.New("invalid prefix not for brasket [" + prefix + "]")
	}
	rightIndex := strings.Index(prefix, "]")
	if rightIndex == -1 {
		return nil, errors.New("invalid prefix for rightIndex [" + prefix + "]")
	}
	name := prefix[1:rightIndex]
	lastName := prefix[rightIndex+1:]
	if name == "" {
		//数组
		curDataArray, ok := curData.([]interface{})
		if !ok {
			curDataArray = []interface{}{}
		}
		singleResult, err := decodeQueryValue(nil, lastName, value)
		if err != nil {
			return nil, err
		}
		curDataArray = append(curDataArray, singleResult)
		return curDataArray, nil
	} else {
		nameInt, err := strconv.Atoi(name)
		if err == nil {
			//数组
			curDataArray, ok := curData.([]interface{})
			if !ok {
				curDataArray = []interface{}{}
			}
			for nameInt >= len(curDataArray) {
				curDataArray = append(curDataArray, "")
			}
			singleResult, err := decodeQueryValue(curDataArray[nameInt], lastName, value)
			if err != nil {
				return nil, err
			}
			curDataArray[nameInt] = singleResult
			return curDataArray, nil
		} else {
			//映射
			curDataMap, ok := curData.(map[string]interface{})
			if !ok {
				curDataMap = map[string]interface{}{}
			}
			singleResult, err := decodeQueryValue(curDataMap[name], lastName, value)
			if err != nil {
				return nil, err
			}
			curDataMap[name] = singleResult
			return curDataMap, nil
		}
	}
}

func decodeQueryInner(data string) (interface{}, error) {
	stringList := strings.Split(data, "&")
	if len(stringList) == 1 && strings.Index(data, "=") == -1 {
		return stringList[0], nil
	}
	resultMap := map[string]interface{}{}
	for _, singleString := range stringList {
		//解析key
		if singleString == "" {
			continue
		}
		singleTemp := strings.Split(singleString, "=")
		var singleKey string
		var singleValue string
		var err error
		if len(singleTemp) == 1 {
			singleKey = singleTemp[0]
			singleValue = ""
		} else {
			singleKey = singleTemp[0]
			singleValue = singleTemp[1]
		}
		if singleKey == "" {
			continue
		}
		singleKey, err = DecodeUrl(singleKey)
		if err != nil {
			return nil, err
		}
		singleValue, err = DecodeUrl(singleValue)
		if err != nil {
			return nil, err
		}
		singleKeyArgv1, singleKeyArgv2 := decodeQueryKey(singleKey)

		//解析value
		curData := resultMap[singleKeyArgv1]
		singleValueResult, err := decodeQueryValue(curData, singleKeyArgv2, singleValue)
		if err != nil {
			return nil, err
		}

		//设置key与value
		resultMap[singleKeyArgv1] = singleValueResult
	}
	return resultMap, nil
}

func DecodeUrlQuery(data []byte, value interface{}) error {
	result, err := decodeQueryInner(string(data))
	if err != nil {
		return err
	}
	return language.MapToArray(result, value, "url")
}
