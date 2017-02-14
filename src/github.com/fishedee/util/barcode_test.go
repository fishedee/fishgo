package util

import (
	"testing"
)

func TestGetQrCodeFile(t *testing.T) {
	testCase := []struct {
		Name string
		Url  string
	}{
		{
			Name: "baidu",
			Url:  "https://www.baidu.com",
		},
		{
			Name: "hongbeibang",
			Url:  "https://www.hongbeibang.com",
		},
	}

	for _, single := range testCase {
		newQrCode := NewQrCode(single.Url, 300, 300)
		singleFile, err := newQrCode.GetQrCodeFile()
		if err != nil {
			panic(err)
		}
		t.Log(singleFile)
	}
}
