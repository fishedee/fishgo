package sdk

import (
	. "github.com/fishedee/util"
	"sync"
	"net"
	"regexp"
	"errors"
)

type Ip138Sdk struct{
}

var mutex sync.Mutex
var currentIP net.IP
var currentIPError error

func (this *Ip138Sdk) GetCurrentIP()(net.IP,error){
	mutex.Lock()
	defer mutex.Unlock()
	if currentIP == nil && currentIPError == nil{
		currentIP,currentIPError = this.getCurrentIPInner()
	}
	return currentIP,currentIPError
}

func (this *Ip138Sdk) getCurrentIPInner()(net.IP,error){
	var result []byte
	var err error

	//获取IP地址
	err = DefaultAjaxPool.Get(&Ajax{
		Url:"http://1111.ip138.com/ic.asp",
		ResponseData:&result,
	})
	if err != nil{
		return nil,err
	}

	//分析IP地址的部分
	reg,err := regexp.Compile("[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+")
	if err != nil{
		return nil,err
	}
	resultIP := reg.Find(result)
	if resultIP == nil{
		return nil,errors.New("缺少IP地址"+string(result))
	}

	//解析IP地址
	return net.ParseIP(string(resultIP)),nil
}
