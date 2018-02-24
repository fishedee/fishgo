package assert

import (
	"testing"
)

func TestDeepEqual(t *testing.T) {
	testCase := []struct {
		left      interface{}
		right     interface{}
		isEqual   bool
		equalDesc string
	}{
		{0, 0, true, ""},
		{nil, nil, true, ""},
		{0, nil, false, "nil != nonil"},
		{int64(0), 0, false, "<int64 Value> type != <int Value> type"},
		{[]string{"123", "456"}, []interface{}{"123", "456"}, false, "<[]string Value> type != <[]interface {} Value> type"},
		{map[string]interface{}{
			"123": "456",
		}, map[string]interface{}{
			"123": 456,
		}, false, "=>123: string type != int type"},
		{map[string]interface{}{
			"123": "456",
		}, map[string]interface{}{
			"123": "789",
		}, false, `=>123: "456" != "789"`},
		{map[string]interface{}{
			"123": "456",
		}, map[string]interface{}{
			"123": "456",
			"456": "123",
		}, false, `: len(map)[1] != len(map)[2]`},
		{map[string]interface{}{
			"123": "456",
		}, map[string]interface{}{
			"456": "123",
		}, false, `=>123: exist != noexist`},
	}
	for index, singleTestCase := range testCase {
		equalDesc, isEqual := DeepEqual(singleTestCase.left, singleTestCase.right)
		if isEqual != singleTestCase.isEqual {
			t.Errorf("%v : [%v] != [%v]", index, isEqual, singleTestCase.isEqual)
		}
		if equalDesc != singleTestCase.equalDesc {
			t.Errorf("%v : [%v] != [%v]", index, equalDesc, singleTestCase.equalDesc)
		}
	}
}
