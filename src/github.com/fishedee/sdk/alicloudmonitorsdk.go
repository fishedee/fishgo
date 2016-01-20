package sdk

import (
	"encoding/json"
	"fmt"
	. "github.com/fishedee/util"
	"sync"
	"time"
)

type AliCloudMonitorSdkData struct {
	metricName string
	value      int64
	timestamp  time.Time
}

type AliCloudMonitorSdk struct {
	AppId     string
	dataMutex sync.Mutex
	data      map[string]AliCloudMonitorSdkData
}

func (this *AliCloudMonitorSdk) getInner(name string) int64 {
	if this.data == nil {
		return 0
	} else {
		value, ok := this.data[name]
		if !ok {
			return 0
		} else {
			return value.value
		}
	}
}

func (this *AliCloudMonitorSdk) setInner(name string, value int64) {
	if this.data == nil {
		this.data = map[string]AliCloudMonitorSdkData{}
	}
	this.data[name] = AliCloudMonitorSdkData{
		metricName: name,
		value:      value,
		timestamp:  time.Now(),
	}
}

func (this *AliCloudMonitorSdk) Max(name string, value int64) {
	this.dataMutex.Lock()
	defer this.dataMutex.Unlock()

	oldvalue := this.getInner(name)
	if oldvalue < value {
		oldvalue = value
	}
	this.setInner(name, oldvalue)
}

func (this *AliCloudMonitorSdk) Min(name string, value int64) {
	this.dataMutex.Lock()
	defer this.dataMutex.Unlock()

	oldvalue := this.getInner(name)
	if oldvalue > value {
		oldvalue = value
	}
	this.setInner(name, oldvalue)
}

func (this *AliCloudMonitorSdk) Set(name string, value int64) {
	this.dataMutex.Lock()
	defer this.dataMutex.Unlock()

	this.setInner(name, value)
}

func (this *AliCloudMonitorSdk) Asc(name string, value int64) {
	this.dataMutex.Lock()
	defer this.dataMutex.Unlock()

	oldvalue := this.getInner(name)
	oldvalue += value
	this.setInner(name, oldvalue)
}

func (this *AliCloudMonitorSdk) Dec(name string, value int64) {
	this.dataMutex.Lock()
	defer this.dataMutex.Unlock()

	oldvalue := this.getInner(name)
	oldvalue -= value
	this.setInner(name, oldvalue)
}

func (this *AliCloudMonitorSdk) Clear() {
	this.dataMutex.Lock()
	defer this.dataMutex.Unlock()

	this.data = map[string]AliCloudMonitorSdkData{}
}

func (this *AliCloudMonitorSdk) GetAllAndClear() map[string]AliCloudMonitorSdkData {
	this.dataMutex.Lock()
	defer this.dataMutex.Unlock()

	result := this.data
	this.data = map[string]AliCloudMonitorSdkData{}
	return result
}

func (this *AliCloudMonitorSdk) Sync() {
	go func() {
		tickChan := time.Tick(time.Minute)
		for {
			<-tickChan
			err := this.syncInner()
			if err != nil {
				fmt.Println("AliCloudMonitorSdk Sync Error " + err.Error())
			}
		}
	}()
}

func (this *AliCloudMonitorSdk) syncInner() error {
	//获取IP
	ip, err := NewIfconfig().GetIP("eth0")
	if err != nil {
		return err
	}
	//获取同步数据
	pushData := this.GetAllAndClear()

	//执行同步
	if pushData == nil || len(pushData) == 0 {
		return nil
	}
	var metrics []interface{}
	for _, singleData := range pushData {
		singleMetric := map[string]interface{}{
			"metricName": singleData.metricName,
			"value":      singleData.value,
			"unit":       "None",
			"dimensions": map[string]string{
				"machineIP": ip.IP.String(),
			},
			"timestamp": singleData.timestamp.Unix()*1000 + int64(singleData.timestamp.Nanosecond()/1000000),
		}
		metrics = append(metrics, singleMetric)
	}
	metricsJson, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	err = DefaultAjaxPool.Get(&Ajax{
		Url: "http://open.cms.aliyun.com/metrics/put",
		Data: map[string]string{
			"userId":    this.AppId,
			"namespace": "acs/custom/" + this.AppId,
			"metrics":   string(metricsJson),
		},
	})
	if err != nil {
		return err
	}
	return nil
}
