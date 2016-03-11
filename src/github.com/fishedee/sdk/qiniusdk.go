package sdk

import (
	"bytes"
	"errors"
	"net/url"

	. "qiniupkg.com/api.v7/kodo"
)

const (
	maxSize = 16 * 1024 * 1024
)

type QiniuSdk struct {
	AccessKey string
	SecretKey string
}

func (this *QiniuSdk) UploadString(bucketName string, data []byte) (string, error) {
	//判断图片文件大小
	fsize := int64(len(data))
	if fsize > maxSize {
		return "", errors.New("上传图片太大！")
	}

	uploadReader := bytes.NewReader(data)
	cfg := &Config{
		AccessKey: this.AccessKey,
		SecretKey: this.SecretKey,
	}
	client := New(0, cfg)
	bucket := client.Bucket(bucketName)
	putRet := PutRet{}
	err := bucket.PutWithoutKey(nil, &putRet, uploadReader, fsize, &PutExtra{})
	if err != nil {
		return "", err
	}

	return putRet.Hash, nil
}

func (this *QiniuSdk) UploadFile(bucketName string, fileAddr string) (string, error) {
	cfg := &Config{
		AccessKey: this.AccessKey,
		SecretKey: this.SecretKey,
	}
	client := New(0, cfg)
	bucket := client.Bucket(bucketName)
	putRet := PutRet{}
	err := bucket.PutFileWithoutKey(nil, &putRet, fileAddr, &PutExtra{})
	if err != nil {
		return "", err
	}

	return putRet.Hash, nil
}

func (this *QiniuSdk) MoveFile(bucketName string, keySrc, keyDest string) error {
	cfg := &Config{
		AccessKey: this.AccessKey,
		SecretKey: this.SecretKey,
	}
	client := New(0, cfg)
	bucket := client.Bucket(bucketName)
	return bucket.Move(nil, keySrc, keyDest)
}

func (this *QiniuSdk) MakeBaseUrl(domain, key string) string {
	return MakeBaseUrl(domain, key)
}

func (this *QiniuSdk) GetUploadToken(bucketName string) (string, error) {
	cfg := &Config{
		AccessKey: this.AccessKey,
		SecretKey: this.SecretKey,
	}
	client := New(0, cfg)

	putPolicy := &PutPolicy{
		Scope:   bucketName,
		Expires: 3600,
	}

	token := client.MakeUptoken(putPolicy)

	return token, nil
}

func (this *QiniuSdk) GetDownloadUrl(inUrl string) (string, error) {
	urlStruct, err := url.Parse(inUrl)
	if err != nil {
		return "", err
	}
	domain := "http://" + urlStruct.Host
	key := urlStruct.Path

	cfg := &Config{
		AccessKey: this.AccessKey,
		SecretKey: this.SecretKey,
	}
	client := New(0, cfg)

	getPolicy := &GetPolicy{
		Expires: 3600,
	}
	baseUrl := MakeBaseUrl(domain, key)
	privateUrl := client.MakePrivateUrl(baseUrl, getPolicy)

	return privateUrl, nil
}
