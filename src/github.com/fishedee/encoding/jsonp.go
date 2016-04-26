package encoding

import (
	"bytes"
	"encoding/json"
	"errors"
	. "github.com/fishedee/language"
	"strings"
)

func EncodeJsonp(functionName string, data interface{}) ([]byte, error) {
	changeValue := ArrayToMap(data, "jsonp")
	jsonResult, err := json.Marshal(changeValue)
	if err != nil {
		return nil, err
	}
	result := []byte(functionName + "(")
	result = append(result, jsonResult...)
	result = append(result, []byte(")")...)
	return result, nil
}

func DecodeJsonp(data []byte, value interface{}) (string, error) {
	leftIndex := bytes.IndexByte(data, '(')
	rightIndex := bytes.LastIndexByte(data, ')')
	if leftIndex == -1 || rightIndex == -1 || leftIndex >= rightIndex {
		return "", errors.New("invalid jsonp format " + string(data))
	}
	functionName := string(data[:leftIndex])
	data = data[leftIndex+1 : rightIndex]
	var valueDynamic interface{}
	err := json.Unmarshal(data, &valueDynamic)
	if err != nil {
		return "", err
	}
	err = MapToArray(valueDynamic, value, "jsonp")
	if err != nil {
		return "", err
	}
	return strings.Trim(functionName, " "), nil
}
