package compress

import (
	"io/ioutil"
	"os"
	"testing"
)

func ClearTestFile(file []string) {
	for _, singleFile := range file {
		os.Remove(singleFile)
	}
}

func TestZipNormal(t *testing.T) {
	testCase := []string{
		"",
		"123",
		"你好",
	}
	txtFile := "zipTest.txt"
	zipFile := "zipTest.zip"
	txtFile2 := "zipTest2.txt"
	defer ClearTestFile([]string{txtFile, zipFile, txtFile2})
	for _, singleTestCase := range testCase {
		var err error
		//准备数据
		ClearTestFile([]string{txtFile, zipFile, txtFile2})
		err = ioutil.WriteFile(txtFile, []byte(singleTestCase), os.ModePerm)
		if err != nil {
			panic(err)
		}
		//压缩
		err = CompressZipFile(txtFile, zipFile)
		if err != nil {
			t.Error("err is not nil! " + err.Error())
			return
		}
		err = DecompressZipFile(zipFile, txtFile2)
		if err != nil {
			t.Error("err is not nil! " + err.Error())
			return
		}
		result, err := ioutil.ReadFile(txtFile2)
		if err != nil {
			t.Error("err is not nil! " + err.Error())
			return
		}
		resultString := string(result)
		if resultString != singleTestCase {
			t.Errorf("%s != %s", resultString, singleTestCase)
			return
		}
	}
}
