package sdk

import (
	"bytes"
	"github.com/upyun/go-sdk/upyun"
)

type UpyunSdk struct {
	Bucket   string
	Operator string
	Password string
	client   *upyun.UpYun
}

func (this *UpyunSdk) getClient() *upyun.UpYun {
	if this.client != nil {
		return this.client
	}
	this.client = upyun.NewUpYun(&upyun.UpYunConfig{
		Bucket:   this.Bucket,
		Operator: this.Operator,
		Password: this.Password,
	})
	return this.client
}

func (this *UpyunSdk) PutString(path string, data []byte) error {
	client := this.getClient()
	return client.Put(&upyun.PutObjectConfig{
		Path:   path,
		Reader: bytes.NewReader(data),
	})
}

func (this *UpyunSdk) PutFile(path string, fileAddr string) error {
	client := this.getClient()
	return client.Put(&upyun.PutObjectConfig{
		Path:      path,
		LocalPath: fileAddr,
	})
}

func (this *UpyunSdk) GetInfo(path string) (*upyun.FileInfo, error) {
	client := this.getClient()
	return client.GetInfo(path)
}

func (this *UpyunSdk) Mkdir(path string) error {
	client := this.getClient()
	return client.Mkdir(path)
}

func (this *UpyunSdk) Delete(path string, isAsync bool) error {
	client := this.getClient()
	return client.Delete(&upyun.DeleteObjectConfig{
		Path:  path,
		Async: isAsync,
	})
}
