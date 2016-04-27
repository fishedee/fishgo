package encoding

import (
	"fmt"
	"reflect"
	"testing"
)

func assertBase64Equal(t *testing.T, left interface{}, right interface{}) {
	isEqual := reflect.DeepEqual(left, right)
	if isEqual == false {
		t.Error(fmt.Sprintf("%#v != %#v", left, right))
	}
}

func TestBase64(t *testing.T) {
	testCase := []struct {
		origin string
		target string
	}{
		{"123", "MTIz"},
		{"123d", "MTIzZA=="},
		{"你好", "5L2g5aW9"},
	}

	for _, singleTestCase := range testCase {
		result, err := EncodeBase64([]byte(singleTestCase.origin))
		assertBase64Equal(t, err, nil)
		assertBase64Equal(t, result, singleTestCase.target)

		result2, err := DecodeBase64(singleTestCase.target)
		assertBase64Equal(t, err, nil)
		assertBase64Equal(t, result2, []byte(singleTestCase.origin))
	}
}
