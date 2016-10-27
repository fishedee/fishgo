package sdk

import (
	"bytes"
	"errors"
	"net/url"

	. "qiniupkg.com/api.v7/kodo"
)

const (
	maxSize = 128 * 1024 * 1024
)

// 七牛sdk
type QiniuSdk struct {
	AccessKey string
	SecretKey string
}

/**
 * [UploadString 上传图片到七牛]
 * @param  string      bucketName [储存区域]
 * @param  []byte      data       [图片字节流]
 * @return string, error          [图片哈希值,错误值]
 */
func (this *QiniuSdk) UploadString(bucketName string, data []byte) (string, error) {
	//判断图片文件大小
	fsize := int64(len(data))
	if fsize > maxSize {
		return "", errors.New("上传图片太大！")
	}

	uploadReader := bytes.NewReader(data)
	bucket := this.getBucket(bucketName)
	putRet := PutRet{}
	err := bucket.PutWithoutKey(nil, &putRet, uploadReader, fsize, &PutExtra{})
	if err != nil {
		return "", err
	}

	return putRet.Hash, nil
}

/**
 * [getBucket 指定储存区域]
 * @param  string     bucketName [储存区域]
 * @return Bucket                [返回储存区域]
 */
func (this *QiniuSdk) getBucket(bucketName string) Bucket {
	client := this.getClient()
	return client.Bucket(bucketName)
}

/**
 * [getClient 连接七牛客户端]
 * @return *Client  [七牛客户端]
 */
func (this *QiniuSdk) getClient() *Client {
	cfg := &Config{
		AccessKey: this.AccessKey,
		SecretKey: this.SecretKey,
	}
	return New(0, cfg)
}

/**
 * [UploadFile 上传文件]
 * @param  string      bucketName [储存区域]
 * @param  string      fileAddr   [本地文件地址]
 * @return string, error          [七牛图片哈希值，错误值]
 */
func (this *QiniuSdk) UploadFile(bucketName string, fileAddr string) (string, error) {
	bucket := this.getBucket(bucketName)
	putRet := PutRet{}
	err := bucket.PutFileWithoutKey(nil, &putRet, fileAddr, &PutExtra{})
	if err != nil {
		return "", err
	}

	return putRet.Hash, nil
}

/**
 * [MoveFile 移动图片]
 * @param  string    bucketName [储存区域]
 * @param  [type]    keySrc     [源路径]
 * @param  string    keyDest    [目的路径]
 * @return error                [错误值]
 */
func (this *QiniuSdk) MoveFile(bucketName string, keySrc, keyDest string) error {
	bucket := this.getBucket(bucketName)
	return bucket.Move(nil, keySrc, keyDest)
}

/**
 * [MakeBaseUrl 获取基本地址]
 * @param  [type]     domain [域名]
 * @param  string     key    [图片哈希值]
 * @return string            [地址]
 */
func (this *QiniuSdk) MakeBaseUrl(domain, key string) string {
	return MakeBaseUrl(domain, key)
}

/**
 * [GetUploadToken 取上传凭证]
 * @param  string      bucketName [储存区域]
 * @return string, error          [凭证，错误值]
 */
func (this *QiniuSdk) GetUploadToken(bucketName string) (string, error) {
	client := this.getClient()

	putPolicy := &PutPolicy{
		Scope:   bucketName,
		Expires: 3600,
	}

	token := client.MakeUptoken(putPolicy)

	return token, nil
}

/**
 * [GetDownloadUrl 取私密下载连接]
 * @param  string      inUrl  [图片链接]
 * @return string, error      [私密链接，错误值]
 */
func (this *QiniuSdk) GetDownloadUrl(inUrl string) (string, error) {
	urlStruct, err := url.Parse(inUrl)
	if err != nil {
		return "", err
	}
	domain := "http://" + urlStruct.Host
	key := urlStruct.Path

	client := this.getClient()

	getPolicy := &GetPolicy{
		Expires: 3600,
	}
	baseUrl := MakeBaseUrl(domain, key)
	privateUrl := client.MakePrivateUrl(baseUrl, getPolicy)

	return privateUrl, nil
}

/**
 * [GetMimeTypeByUrl 取文件的MIME类型]
 * @param  [type]      bucketName [储存区域]
 * @param  string      url        [图片哈希]
 * @return string, error          [类型，错误值]
 */
func (this *QiniuSdk) GetMimeTypeByKey(bucketName, key string) (string, error) {
	bucket := this.getBucket(bucketName)
	stat, err := bucket.Stat(nil, key)
	if err != nil {
		return "", err
	}
	return stat.MimeType, nil
}
