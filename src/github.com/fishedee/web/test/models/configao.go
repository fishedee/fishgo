package test

import (
	. "github.com/fishedee/web"
	"time"
)

type configAoModel struct {
	Model
	dataStruct map[string]ConfigData
}

type ConfigData struct {
	Data       string
	CreateTime time.Time
	ModifyTime time.Time
}

func (this *configAoModel) Set(key string, value string) {
	this.Cache.Set(key, value, time.Hour*24)
}

func (this *configAoModel) Get(key string) string {
	data, _ := this.Cache.Get(key)
	return data
}

func (this *configAoModel) SetStruct(key string, value ConfigData) {
	if this.dataStruct == nil {
		this.dataStruct = map[string]ConfigData{}
	}
	value.CreateTime = time.Now()
	value.ModifyTime = time.Now()
	this.dataStruct[key] = value
}

func (this *configAoModel) GetStruct(key string) ConfigData {
	return this.dataStruct[key]
}
