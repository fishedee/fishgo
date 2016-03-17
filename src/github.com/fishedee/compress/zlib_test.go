package compress

import (
	"testing"
)

func TestZlibNormal(t *testing.T) {
	testCase := []string{
		"",
		"123",
		"你好",
	}
	for _, singleTestCase := range testCase {
		//压缩
		data, err := CompressZlib([]byte(singleTestCase))
		if err != nil {
			t.Error("err is not nil! " + err.Error())
			return
		}
		dataResult, err := DecompressZlib(data)
		if err != nil {
			t.Error("err is not nil! " + err.Error())
			return
		}
		stringResult := string(dataResult)
		if stringResult != singleTestCase {
			t.Errorf("%s != %s", stringResult, singleTestCase)
			return
		}
	}
}
