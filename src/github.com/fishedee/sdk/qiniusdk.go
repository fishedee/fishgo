package sdk

import (
	"io"
	"net/url"
	. "qiniupkg.com/api.v7/kodo"
)

type QiniuSdk struct {
	AccessKey string
	SecretKey string
}

func (this *QiniuSdk) UploadString(bucketName string, data io.Reader, fsize int64) (string, error) {
	cfg := &Config{
		AccessKey: this.AccessKey,
		SecretKey: this.SecretKey,
	}
	client := New(0, cfg)
	bucket := client.Bucket(bucketName)
	putRet := PutRet{}
	err := bucket.PutWithoutKey(nil, &putRet, data, fsize, &PutExtra{})
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
