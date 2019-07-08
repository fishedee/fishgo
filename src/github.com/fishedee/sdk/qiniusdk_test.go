// 七牛文件上传接口测试
package sdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

// 测试七牛字节流上传接口--断点续传
func TestRuploadString(t *testing.T) {
	testRuploadString(t)
}

func testRuploadString(t *testing.T) {
	testCase := getTestCase()

	qiniuSdk := newQiniuSdk()
	for _, single := range testCase {
		content, err := ioutil.ReadFile(single.file)
		if err != nil {
			panic(err)
		}
		t.Log(getUseTime(func() { qiniuSdk.RuploadString(single.bucket, content) }))
	}
}

// 测试七牛文件上传接口--断点续传
func TestRuploadFile(t *testing.T) {
	testRuploadFile(t)
}

func testRuploadFile(t *testing.T) {
	testCase := getTestCase()

	qiniuSdk := newQiniuSdk()
	for _, single := range testCase {
		t.Log(getUseTime(func() { qiniuSdk.RuploadFile(single.bucket, single.file) }))
	}
}

// 测试七牛字节流上传接口
func TestUploadString(t *testing.T) {
	testUploadString(t)
}

func testUploadString(t *testing.T) {
	testCase := getTestCase()

	qiniuSdk := newQiniuSdk()
	for _, single := range testCase {
		content, err := ioutil.ReadFile(single.file)
		if err != nil {
			panic(err)
		}
		t.Log(getUseTime(func() { qiniuSdk.UploadString(single.bucket, content) }))
	}
}

// 测试七牛文件上传接口
func TestUploadFile(t *testing.T) {
	testUploadFile(t)
}

func testUploadFile(t *testing.T) {
	testCase := getTestCase()

	qiniuSdk := newQiniuSdk()
	for _, single := range testCase {
		t.Log(getUseTime(func() { qiniuSdk.UploadFile(single.bucket, single.file) }))
	}
}

// 测试用例
type uploadTestCase struct {
	name   string
	file   string
	bucket string
}

// 初始化测试用例
func getTestCase() []uploadTestCase {
	return []uploadTestCase{
		{"localimage", "./testdata/test.jpg", "bakeweb"},
		{"localaudio", "./testdata/test.aac", "bakewebaudio"},
		{"localvideo", "./testdata/test.mp4", "bakewebvideo"},
	}
}

// 建立测试用实例
func newQiniuSdk() QiniuSdk {
	config, err := ioutil.ReadFile("./testdata/config.json")
	if err != nil {
		panic(err)
	}
	var data struct {
		Qiniu struct {
			AccessKey string
			SecretKey string
		}
	}
	err = json.Unmarshal(config, &data)
	if err != nil {
		panic(err)
	}
	return QiniuSdk{
		AccessKey: data.Qiniu.AccessKey,
		SecretKey: data.Qiniu.SecretKey,
	}
}

// 获取操作时间
func getUseTime(f func()) int64 {
	beginTime := time.Now()
	beginTimeNano := beginTime.UnixNano()
	f()
	endTime := time.Now()
	endTimeNano := endTime.UnixNano()
	fmt.Println(beginTime, endTime)
	return int64(endTimeNano - beginTimeNano)
}
