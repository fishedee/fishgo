package test

import (
	. "github.com/fishedee/web"
	"time"
)

type ConfigAoModel struct {
	Model
	dataStruct map[string]ConfigData
}

type ConfigData struct {
	Data       string
	CreateTime time.Time
	ModifyTime time.Time
}

func (this *ConfigAoModel) Set(key string, value string) {
	this.Cache.Set(key, value, time.Hour*24)
}

func (this *ConfigAoModel) Get(key string) string {
	data, _ := this.Cache.Get(key)
	return data
}

func (this *ConfigAoModel) SetStruct(key string, value ConfigData) {
	if this.dataStruct == nil {
		this.dataStruct = map[string]ConfigData{}
	}
	value.CreateTime = time.Now()
	value.ModifyTime = time.Now()
	this.dataStruct[key] = value
}

func (this *ConfigAoModel) GetStruct(key string) ConfigData {
	return this.dataStruct[key]
}
