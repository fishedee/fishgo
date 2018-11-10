package sdk

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pili-engineering/pili-sdk-go.v2/pili"
	"io"
	"qiniupkg.com/api.v7/kodo"
)

const (
	maxSize = 256 * 1024 * 1024 // 文件最大容量
)

// 七牛sdk
type QiniuSdk struct {
	AccessKey string
	SecretKey string
	Zone      int
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
	putRet := kodo.PutRet{}
	err := bucket.PutWithoutKey(nil, &putRet, uploadReader, fsize, &kodo.PutExtra{})
	if err != nil {
		return "", err
	}

	return putRet.Hash, nil
}

/**
 * [UploadString 上传文件到七牛--支持断点续传]
 * @param  string      bucketName [储存区域]
 * @param  []byte      data       [图片字节流]
 * @return string, error          [图片哈希值,错误值]
 */
func (this *QiniuSdk) RuploadString(bucketName string, data []byte) (string, error) {
	//判断图片文件大小
	fsize := int64(len(data))
	if fsize > maxSize {
		return "", errors.New("上传图片太大！")
	}

	ctx := context.WithValue(context.Background(), 1, 1)
	uploadReader := bytes.NewReader(data)
	bucket := this.getBucket(bucketName)
	putRet := kodo.PutRet{}
	putExtra := kodo.RputExtra{
		ChunkSize: 4 * 1024 * 1024, // 每段大小
		TryTimes:  5,               // 尝试次数
	}
	err := bucket.RputWithoutKey(ctx, &putRet, uploadReader, fsize, &putExtra)
	if err != nil {
		return "", err
	}

	return putRet.Hash, nil
}

/**
 * [UploadFile 上传文件]
 * @param  string      bucketName [储存区域]
 * @param  string      fileAddr   [本地文件地址]
 * @return string, error          [七牛图片哈希值，错误值]
 */
func (this *QiniuSdk) UploadFile(bucketName, fileAddr string) (string, error) {
	// 检查文件大小
	err := this.checkFileSize(fileAddr)
	if err != nil {
		return "", err
	}

	bucket := this.getBucket(bucketName)
	putRet := kodo.PutRet{}
	err = bucket.PutFileWithoutKey(nil, &putRet, fileAddr, &kodo.PutExtra{})
	if err != nil {
		return "", err
	}

	return putRet.Hash, nil
}

/**
 * [UploadString 上传文件到七牛--支持断点续传]
 * @param  string      bucketName [储存区域]
 * @param  []byte      data       [图片字节流]
 * @return string, error          [图片哈希值,错误值]
 */
func (this *QiniuSdk) RuploadFile(bucketName, fileAddr string) (string, error) {
	// 检查文件大小
	err := this.checkFileSize(fileAddr)
	if err != nil {
		return "", err
	}

	// 上传
	ctx := context.WithValue(context.Background(), 1, 1)
	bucket := this.getBucket(bucketName)
	putRet := kodo.PutRet{}
	putExtra := kodo.RputExtra{
		ChunkSize: 4 * 1024 * 1024, // 每段大小
		TryTimes:  5,               // 尝试次数
	}
	err = bucket.RputFileWithoutKey(ctx, &putRet, fileAddr, &putExtra)
	if err != nil {
		return "", err
	}

	return putRet.Hash, nil
}

/**
* [List 列举七牛的文件]
* @param  string      bucketName [储存区域]
* @param  string      prefix     [指定前缀，只有资源名匹配该前缀的资源会被列出]
* @param  string      delimiter  [指定目录分隔符，列出所有公共前缀（模拟列出目录效果）]
* @param  string      marker     [上一次列举返回的位置标记，作为本次列举的起点信息]
* @param  string      limit      [本次列举的条目数，范围为1-1000]
* @return kodo.ListItem,[]string,string,error          [返回条目的数组,返回目录名的数组,有剩余条目则返回非空字符串，作为下一次列举的参数传入]
 */
func (this *QiniuSdk) List(bucketName string, prefix, delimiter, marker string, limit int) ([]kodo.ListItem, []string, string, error) {
	bucket := this.getBucket(bucketName)
	entries, commonPrefixes, markerOut, err := bucket.List(nil, prefix, delimiter, marker, limit)
	if err == io.EOF {
		err = nil
	}
	return entries, commonPrefixes, markerOut, err
}

// 检查文件大小
func (this *QiniuSdk) checkFileSize(fileAddr string) error {
	fileInfo, err := os.Stat(fileAddr)
	if err != nil {
		return err
	}
	if fileInfo.Size() > maxSize {
		return errors.New("上传图片太大！")
	}
	return nil
}

/**
 * [getBucket 指定储存区域]
 * @param  string     bucketName [储存区域]
 * @return Bucket                [返回储存区域]
 */
func (this *QiniuSdk) getBucket(bucketName string) kodo.Bucket {
	client := this.getClient()
	return client.Bucket(bucketName)
}

/**
 * [getClient 连接七牛客户端]
 * @return *Client  [七牛客户端]
 */
func (this *QiniuSdk) getClient() *kodo.Client {
	cfg := &kodo.Config{
		AccessKey: this.AccessKey,
		SecretKey: this.SecretKey,
	}
	zone := this.Zone
	if zone == 2 {
		//华南地区
		zone = 0
		cfg.UpHosts = []string{"http://up-z2.qiniup.com", "http://upload-z2.qiniup.com"}
		cfg.IoHost = "http://iovip-z2.qbox.me"
		cfg.RSHost = "http://rs-z2.qbox.me"
		cfg.RSFHost = "http://rsf-z2.qbox.me"
	}
	return kodo.New(zone, cfg)
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
	return kodo.MakeBaseUrl(domain, key)
}

/**
 * [GetUploadToken 取上传凭证]
 * @param  string      bucketName [储存区域]
 * @return string, error          [凭证，错误值]
 */
func (this *QiniuSdk) GetUploadToken(bucketName string) (string, error) {
	client := this.getClient()

	putPolicy := &kodo.PutPolicy{
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
	domain := urlStruct.Host
	key := strings.Replace(urlStruct.Path, "/", "", -1)

	client := this.getClient()

	getPolicy := &kodo.GetPolicy{
		Expires: 86400,
	}
	baseUrl := kodo.MakeBaseUrl(domain, key)
	privateUrl := client.MakePrivateUrl(baseUrl, getPolicy)

	return privateUrl, nil
}

/**
 * [GetDownloadUrl 取私密视频缩略图下载连接]
 * @param  string      inUrl  [图片链接]
 * @param  string      sample  [缩略图生产参数]
 * @return string, error      [私密链接，错误值]
 */
func (this *QiniuSdk) GetSampleDownloadUrl(inUrl, sample string) (string, error) {
	urlStruct, err := url.Parse(inUrl)
	if err != nil {
		return "", err
	}
	domain := urlStruct.Host
	key := strings.Replace(urlStruct.Path, "/", "", -1)

	client := this.getClient()

	getPolicy := &kodo.GetPolicy{
		Expires: 3600,
	}
	baseUrl := kodo.MakeBaseUrl(domain, key)
	privateUrl := client.MakePrivateUrl(baseUrl+"?"+sample, getPolicy)

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

/**
 * [NewMac 授权信息]
 * @return *MAC  [授权信息]
 */
func (this *QiniuSdk) NewMac() *pili.MAC {
	return &pili.MAC{
		AccessKey: this.AccessKey,
		SecretKey: []byte(this.SecretKey),
	}
}

/**
 * [GetRTMPPublishURL 生成 RTMP 推流地址]
 * @param  [type]     domain             [与直播空间绑定的 RTMP 推流域名，可以在 portal.qiniu.com 上绑定]
 * @param  [type]     hubName                [直播空间名称]
 * @param  string     streamKey          [流名，流不需要事先存在，推流会自动创建流]
 * @param  *MAC       mac                [授权信息]
 * @param  int64      expireAfterSeconds [生成的推流地址的有效时间]
 * @return string                                 [RTMP推流地址]
 */
func (this *QiniuSdk) GetRTMPPublishURL(domain, hubName, streamKey string, mac *pili.MAC, expireAfterSeconds int64) string {
	return pili.RTMPPublishURL(domain, hubName, streamKey, mac, expireAfterSeconds)
}

/**
 * [GetRTMPPlayURL 生成 RTMP 播放地址]
 * @param  [type]     domain    [绑定的直播域名]
 * @param  [type]     hubName       [直播空间名称]
 * @param  string     streamKey [流名，流不需要事先存在，推流会自动创建流]
 * @return string                        [播放地址]
 */
func (this *QiniuSdk) GetRTMPPlayURL(domain, hubName, streamKey string) string {
	return pili.RTMPPlayURL(domain, hubName, streamKey)
}

/**
 * [GetSnapshotPlayURL 生成直播封面地址]
 * @param  [type]     domain    [绑定的直播域名]
 * @param  [type]     hubName       [直播空间名称]
 * @param  string     streamKey [流名，流不需要事先存在，推流会自动创建流]
 * @return string                        [生成直播封面地址]
 */
func (this *QiniuSdk) GetSnapshotPlayURL(domain, hubName, streamKey string) string {
	return pili.SnapshotPlayURL(domain, hubName, streamKey)
}

/**
 * [NewClient 初始化授权客户]
 * @param  *pili.MAC         mac [授权信息]
 * @param  http.RoundTripper tr  [http会话]
 * @return *pili.Client                   [授权客户]
 */
func (this *QiniuSdk) NewClient(mac *pili.MAC, tr http.RoundTripper) *pili.Client {
	return pili.New(mac, tr)
}
