package util

import (
	"testing"
)

func TestSnapshot(t *testing.T) {
	testCase := []struct {
		Name string // 用例名称
		Url  string // 文件路径
		Path string // 生成图片保存文件
	}{
		{
			Name: "hongbeibang",
			Url:  "http://www.hongbeibang.com",
			Path: "testdata/a.png",
		},
		{
			Name: "local file",
			Url:  "testdata/localfile.html",
			Path: "testdata/b.png",
		},
	}

	// userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.94 Safari/537.36"
	userAgent := "Mozilla/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X; en-us) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11A465 Safari/9537.53"
	width := 400
	height := 640
	for _, single := range testCase {
		brw := NewBrowser(userAgent, "", 0, width, height)
		err := brw.Snapshot(single.Url, single.Path)
		if err != nil {
			panic(err)
		}
	}
}
