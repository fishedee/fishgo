package util

import (
	"io/ioutil"
	"strings"
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
		{
			Name: "invite file",
			Url:  "testdata/invitation.html",
			Path: "testdata/c.png",
		},
	}

	// userAgent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/37.0.2062.94 Safari/537.36"
	userAgent := "Mozilla/5.0 (iPhone; CPU iPhone OS 7_0 like Mac OS X; en-us) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11A465 Safari/9537.53"
	width := 400
	height := 640
	for _, single := range testCase {
		brw := NewBrowser(userAgent, "", 0, width, height)
		if single.Name == "invite file" {
			fileData, err := ioutil.ReadFile(single.Url)
			if err != nil {
				panic(err)
			}
			fileDataStr := string(fileData)
			fileDataStr = strings.Replace(fileDataStr, "{{userImage}}", "http://image.hongbeibang.com/FpsiNicGwECZnreITGMB1AEJVriR?750X561", -1)
			fileDataStr = strings.Replace(fileDataStr, "{{userName}}", "jd", -1)
			fileDataStr = strings.Replace(fileDataStr, "{{educationCourseTitle}}", "戚风", -1)
			fileDataStr = strings.Replace(fileDataStr, "{{courseDegreeTitle}}", "本科", -1)
			fileDataStr = strings.Replace(fileDataStr, "{{courseSummary}}", "戚风入门来看我，包你一学就会。", -1)
			fileDataStr = strings.Replace(fileDataStr, "{{courseBeginTime}}", "2017-02-20 19:00", -1)
			fileDataStr = strings.Replace(fileDataStr, "{{courseQrCode}}", "http://image.hongbeibang.com/FjAksoRCYDdBiTb3fld6p0v7n-7Z?300X300", -1)
			fileName, err := CreateTempFile("invite", ".html")
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile(fileName, []byte(fileDataStr), 0666)
			if err != nil {
				panic(err)
			}
			t.Log(fileName)
			single.Url = fileName
		}
		err := brw.Snapshot(single.Url, single.Path)
		if err != nil {
			panic(err)
		}
	}
}
