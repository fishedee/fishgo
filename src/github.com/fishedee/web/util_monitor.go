package web

import (
	"errors"
	. "github.com/fishedee/sdk"
)

type Monitor interface {
	AscErrorCount()
	AscCriticalCount()
}

type MonitorConfig struct {
	Driver        string
	AppId         string
	ErrorCount    string
	CriticalCount string
}

type monitorImplement struct {
	AliCloudMonitorSdk
	config MonitorConfig
}

func NewMonitor(config MonitorConfig) (Monitor, error) {
	if config.Driver == "" {
		return nil, nil
	} else if config.Driver == "aliyuncloudmonitor" {
		result := &monitorImplement{
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

func NewMonitorFromConfig(configName string) (Monitor, error) {
	driver := globalBasic.Config.GetString(configName + "driver")
	appId := globalBasic.Config.GetString(configName + "appid")
	errorCount := globalBasic.Config.GetString(configName + "errorcount")
	criticalCount := globalBasic.Config.GetString(configName + "criticalcount")

	monitorConfig := MonitorConfig{}
	monitorConfig.Driver = driver
	monitorConfig.AppId = appId
	monitorConfig.ErrorCount = errorCount
	monitorConfig.CriticalCount = criticalCount
	return NewMonitor(monitorConfig)
}

func (this *monitorImplement) AscErrorCount() {
	if this.config.ErrorCount != "" {
		this.Asc(this.config.ErrorCount, 1)
	}
}

func (this *monitorImplement) AscCriticalCount() {
	if this.config.CriticalCount != "" {
		this.Asc(this.config.CriticalCount, 1)
	}
}
