package util

import (
	"os"
	"strings"
	"testing"
)

func TestTemp(t *testing.T) {
	testCase := []struct {
		dir    string
		suffix string
	}{
		{"", ""},
		{"", ".go"},
		{"mc", ""},
		{"mc", ".go"},
	}
	for _, singleTestCase := range testCase {
		tempFile, err := CreateTempFile(singleTestCase.dir, singleTestCase.suffix)
		if err != nil {
			t.Errorf("tempFile has error %s", err.Error())
			return
		}
		if singleTestCase.dir != "" {
			if strings.HasPrefix(tempFile, strings.TrimRight(os.TempDir(), "/")+"/"+singleTestCase.dir) == false {
				t.Errorf("tempFile has no dir %s", singleTestCase.dir)
				return
			}
		}
		if singleTestCase.suffix != "" {
			if strings.HasSuffix(tempFile, singleTestCase.suffix) == false {
				t.Errorf("tempFile has no suffix %s", singleTestCase.suffix)
				return
			}
		}
	}
}
