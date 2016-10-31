package sdk

import (
	"strconv"
	"strings"
	"time"

	"github.com/fishedee/crypto"
	"github.com/fishedee/encoding"
	"github.com/fishedee/util"
)

type NetEaseSdk struct {
	AppKey    string
	AppSecret string
}

type NetEaseCommonParam struct {
	AppKey   string //
	Nonce    string // 随机数，最大长度128个字符
	CurTime  string // 时间戳,秒数
	CheckSum string // sha1(AppSecret+Nonce+CurTime),转16进制小写
}

type NetEaseParam struct {
	Type       int
	Records    int
	Pnum       int
	Sort       int
	NeedRecord int
	Format     int
	Duration   int
	Cid        string
	Name       string
	Ofield     string
	Filename   string
	CidList    []string
}

type NetEaseCommonRet struct {
	Code int
	Msg  string
}

type NetEaseCreateChanRet struct {
	Cid         string // 频道iD，32位字符串
	Ctime       int    // 创建频道的时间戳
	Name        string // 频道名称
	PushUrl     string // 推流地址
	HttpPullUrl string // http拉流地址
	HlsPullUrl  string // hls拉流地址
	RtmpPullUrl string // rtmp拉流地址
}

type NetEaseCreateChanResult struct {
	NetEaseCommonRet
	Ret NetEaseCreateChanRet
}

/**
 * [CreateChan 创建频道]
 * @param  string    name [频道名字]
 * @param  int       type [频道类型：0为rtmp]
 * @return NetEaseCreateChanResult, error                   [频道信息，错误值]
 */
func (this *NetEaseSdk) CreateChan(name string, ctype int) (NetEaseCreateChanResult, error) {
	url := "https://vcloud.163.com/app/channel/create"
	method := "post"
	result := NetEaseCreateChanResult{}
	data := NetEaseParam{
		Name: name,
		Type: ctype,
	}
	err := this.api(url, method, data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

/**
 * [UpdateChan 更新频道]
 * @param  string             name   [频道名称]
 * @param  string             cid    [频道ID]
 * @param  int                ctype  [频道类型]
 * @return NetEaseCommonRet, error               [修改结果，错误值]
 */
func (this *NetEaseSdk) UpdateChan(name string, cid string, ctype int) (NetEaseCommonRet, error) {
	url := "https://vcloud.163.com/app/channel/update"
	method := "post"
	result := NetEaseCommonRet{}
	data := NetEaseParam{
		Name: name,
		Cid:  cid,
		Type: ctype,
	}
	err := this.api(url, method, data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

/**
 * [DeleteChan 删除频道]
 * @param  string             cid    [频道ID]
 * @return NetEaseCommonRet, error               [结果,错误值]
 */
func (this *NetEaseSdk) DeleteChan(cid string) (NetEaseCommonRet, error) {
	url := "https://vcloud.163.com/app/channel/delete"
	method := "post"
	result := NetEaseCommonRet{}
	data := NetEaseParam{
		Cid: cid,
	}
	err := this.api(url, method, data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

type NetEaseStatRet struct {
	Ctime      int    // 创建频道的时间戳
	Cid        string // 频道ID，32位字符串
	Name       string // 频道名称
	Status     int    // 频道状态（0：空闲； 1：直播； 2：禁用； 3：直播录制）
	Type       int    // 频道类型 ( 0 : rtmp, 1 : hls, 2 : http)
	Uid        int    // 用户ID
	NeedRecord int    // 1-开启录制； 0-关闭录制
	Format     int    // 1-flv； 0-mp4
	Duration   int    // 录制切片时长(分钟)，默认120分钟
	Filename   string // 录制后文件名
	OnlineUser int    // 在线用户
}

type NetEaseStatResult struct {
	NetEaseCommonRet
	Ret NetEaseStatRet
}

/**
 * [ChanStat 获取频道状态]
 * @param  string                 cid    [频道ID]
 * @return NetEaseStatResult, error               [状态结果，错误值]
 */
func (this *NetEaseSdk) ChanStat(cid string) (NetEaseStatResult, error) {
	url := "https://vcloud.163.com/app/channelstats"
	method := "post"
	result := NetEaseStatResult{}
	data := NetEaseParam{
		Cid: cid,
	}
	err := this.api(url, method, data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

type NetEaseListRet struct {
	Pnum int
	List []NetEaseStatRet
}

type NetEaseListResult struct {
	NetEaseCommonRet
	Ret NetEaseListRet
}

/**
 * [ChanList 频道列表]
 * @param  [type]                 records [单页记录数，默认值为10]
 * @param  int                    pnum    [要取第几页，默认值为1]
 * @param  string                 ofield  [排序的域，支持的排序域为：ctime（默认）]
 * @param  int                    sort    [升序还是降序，1升序，0降序，默认为desc]
 * @return NetEaseListResult, error                [频道列表，错误值]
 */
func (this *NetEaseSdk) ChanList(records, pnum int, ofield string, sort int) (NetEaseListResult, error) {
	url := "https://vcloud.163.com/app/channellist"
	method := "post"
	result := NetEaseListResult{}
	data := NetEaseParam{
		Records: records,
		Pnum:    pnum,
		Ofield:  ofield,
		Sort:    sort,
	}
	err := this.api(url, method, data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

/**
 * [ChanAddr 重新获得频道地址]
 * @param  string                       cid    [频道ID]
 * @return NetEaseCreateChanResult, error               [频道地址信息,错误值]
 */
func (this *NetEaseSdk) ChanAddr(cid string) (NetEaseCreateChanResult, error) {
	url := "https://vcloud.163.com/app/address"
	method := "post"
	result := NetEaseCreateChanResult{}
	data := NetEaseParam{
		Cid: cid,
	}
	err := this.api(url, method, data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

/**
 * [SetChanRecord 设置频道为录制状态]
 * @param  string                cid        [频道ID，32位字符串]
 * @param  [type]                needRecord [1-开启录制； 0-关闭录制]
 * @param  [type]                format     [1-flv； 0-mp4]
 * @param  int                   duration   [录制切片时长(分钟)，5~120分钟]
 * @param  string                filename   [录制后文件名（只支持中文、字母和数字），格式为filename_YYYYMMDD-HHmmssYYYYMMDD-HHmmss, 文件名录制起始时间（年月日时分秒) -录制结束时间（年月日时分秒)]
 * @return NetEaseCommonRet, error                   [状态，错误值]
 */
func (this *NetEaseSdk) SetChanRecord(cid string, needRecord, format, duration int, filename string) (NetEaseCommonRet, error) {
	url := "https://vcloud.163.com/app/channel/setAlwaysRecord"
	method := "post"
	result := NetEaseCommonRet{}
	data := NetEaseParam{
		Cid:        cid,
		NeedRecord: needRecord,
		Format:     format,
		Duration:   duration,
		Filename:   filename,
	}
	err := this.api(url, method, data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

/**
 * [PauseChan 禁用用户正在直播的频道]
 * @param  string                cid    [频道ID，32位字符串]
 * @return NetEaseCommonRet, error               [状态，错误值]
 */
func (this *NetEaseSdk) PauseChan(cid string) (NetEaseCommonRet, error) {
	url := "https://vcloud.163.com/app/channel/pause"
	method := "post"
	result := NetEaseCommonRet{}
	data := NetEaseParam{
		Cid: cid,
	}
	err := this.api(url, method, data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

type NetEasePauseListRet struct {
	SuccessList []string
}
type NetEasePauseListResult struct {
	NetEaseCommonRet
	Ret NetEasePauseListRet
}

/**
 * [PauseChanList 禁用一组用户正在直播的频道]
 * @param  []string              cids   [频道ID列表]
 * @return NetEasePauseListResult, error               [description]
 */
func (this *NetEaseSdk) PauseChanList(cids []string) (NetEasePauseListResult, error) {
	url := "https://vcloud.163.com/app/channellist/pause"
	method := "post"
	result := NetEasePauseListResult{}
	data := NetEaseParam{
		CidList: cids,
	}
	err := this.api(url, method, data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

/**
 * [ResumeChan 恢复用户被禁用的频道]
 * @param  string                cid    [频道ID，32位字符串]
 * @return NetEaseCommonRet, error               [description]
 */
func (this *NetEaseSdk) ResumeChan(cid string) (NetEaseCommonRet, error) {
	url := "https://vcloud.163.com/app/channel/resume"
	method := "post"
	result := NetEaseCommonRet{}
	data := NetEaseParam{
		Cid: cid,
	}
	err := this.api(url, method, data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

/**
 * [ResumeChanList 恢复一组用户正在直播的频道]
 * @param  []string                    cids   [频道ID列表]
 * @return NetEasePauseListResult, error               [description]
 */
func (this *NetEaseSdk) ResumeChanList(cids []string) (NetEasePauseListResult, error) {
	url := "https://vcloud.163.com/app/channellist/resume"
	method := "post"
	result := NetEasePauseListResult{}
	data := NetEaseParam{
		CidList: cids,
	}
	err := this.api(url, method, data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

type NetEaseVideoRet struct {
	VideoName    string `json:"video_name"`
	OrigVideoKey string `json:"orig_video_key"`
	Uid          int
	vid          int
}

type NetEaseVideoListRet struct {
	VideoList []NetEaseVideoRet
}

type NetEaseVideoListResult struct {
	NetEaseCommonRet
	Ret NetEaseVideoListRet
}

/**
 * [VideoList 获取某频道录制视频文件列表]
 * @param  string                      cid    [频道ID，32位字符串]
 * @return NetEaseVideoListResult, error               [description]
 */
func (this *NetEaseSdk) VideoList(cid string) (NetEaseVideoListResult, error) {
	url := "https://vcloud.163.com/app/videolist"
	method := "post"
	result := NetEaseVideoListResult{}
	data := NetEaseParam{
		Cid: cid,
	}
	err := this.api(url, method, data, &result)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (this *NetEaseSdk) api(url, method string, data interface{}, responseData interface{}) error {
	var dataByte []byte
	nonce := crypto.CryptoRand(64)
	curTime := strconv.FormatInt(time.Now().Unix(), 10)
	err := util.DefaultAjaxPool.Post(&util.Ajax{
		Url:    url,
		Data:   data,
		Method: method,
		Header: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
			"AppKey":       this.AppKey,
			"Nonce":        nonce,
			"CurTime":      curTime,
			"CheckSum":     this.getSignature(nonce, curTime),
		},
		ResponseData: &dataByte,
	})
	if err != nil {
		return err
	}
	err = encoding.DecodeJson(dataByte, &responseData)
	if err != nil {
		return err
	}
	return nil
}

func (this *NetEaseSdk) getSignature(nonce, curTime string) string {
	sha1Str := crypto.CryptoSha1([]byte(this.AppSecret + nonce + curTime))
	return strings.ToLower(sha1Str)
}
