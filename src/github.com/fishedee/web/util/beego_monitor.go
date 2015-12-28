package util

import (
	"errors"
	"github.com/astaxie/beego"
	. "github.com/fishedee/sdk"
	. "github.com/fishedee/util"
)

type MonitorManagerConfig struct {
	Driver        string
	AppId         string
	ErrorCount    string
	CriticalCount string
}

type MonitorManager struct {
	AliCloudMonitorSdk
	config MonitorManagerConfig
}

var newMonitorManagerMemory *MemoryFunc
var newMonitorManagerFromConfigMemory *MemoryFunc

func init() {
	var err error
	newMonitorManagerMemory, err = NewMemoryFunc(newMonitorManager, MemoryFuncCacheNormal)
	if err != nil {
		panic(err)
	}
	newMonitorManagerFromConfigMemory, err = NewMemoryFunc(newMonitorManagerFromConfig, MemoryFuncCacheNormal)
	if err != nil {
		panic(err)
	}
}

func newMonitorManager(config MonitorManagerConfig) (*MonitorManager, error) {
	if config.Driver == "" {
		return nil, nil
	} else if config.Driver == "aliyuncloudmonitor" {
		result := &MonitorManager{
			AliCloudMonitorSdk: AliCloudMonitorSdk{
				AppId: config.AppId,
			},
			config: config,
		}
		go result.AliCloudMonitorSdk.Sync()
		return result, nil
	} else {
		return nil, errors.New("invalid monitor config " + config.Driver)
	}
}

func NewMonitorManager(config MonitorManagerConfig) (*MonitorManager, error) {
	result, err := newMonitorManagerMemory.Call(config)
	if err != nil {
		return nil, err
	}
	return result.(*MonitorManager), err
}

func newMonitorManagerFromConfig(configName string) (*MonitorManager, error) {
	driver := beego.AppConfig.String(configName + "driver")
	appId := beego.AppConfig.String(configName + "appid")
	errorCount := beego.AppConfig.String(configName + "errorcount")
	criticalCount := beego.AppConfig.String(configName + "criticalcount")

	monitorConfig := MonitorManagerConfig{}
	monitorConfig.Driver = driver
	monitorConfig.AppId = appId
	monitorConfig.ErrorCount = errorCount
	monitorConfig.CriticalCount = criticalCount
	return NewMonitorManager(monitorConfig)
}

func NewMonitorManagerFromConfig(configName string) (*MonitorManager, error) {
	result, err := newMonitorManagerFromConfigMemory.Call(configName)
	if err != nil {
		return nil, err
	}
	return result.(*MonitorManager), err
}

func (this *MonitorManager) AscErrorCount() {
	if this.config.ErrorCount != "" {
		this.Asc(this.config.ErrorCount, 1)
	}
}

func (this *MonitorManager) AscCriticalCount() {
	if this.config.CriticalCount != "" {
		this.Asc(this.config.CriticalCount, 1)
	}
}
