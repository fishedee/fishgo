package compress

import (
	"testing"
)

func TestGzipNormal(t *testing.T) {
	testCase := []string{
		"",
		"123",
		"你好",
	}
	for _, singleTestCase := range testCase {
		//压缩
		data, err := CompressGzip([]byte(singleTestCase))
		if err != nil {
			t.Error("err is not nil! " + err.Error())
			return
		}
		dataResult, err := DecompressGzip(data)
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
