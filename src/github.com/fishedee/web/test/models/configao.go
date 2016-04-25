package test

import (
	"time"
)

type ConfigAoModel struct {
	BaseModel
	dataStruct map[string]ConfigData
}

type ConfigData struct {
	Data       string
	CreateTime time.Time `xorm:"created"`
	ModifyTime time.Time `xorm:"updated"`
}

func (this *ConfigAoModel) Set(key string, value string) {
	this.Cache.Set(key, value, time.Hour*24)
}

func (this *ConfigAoModel) Get(key string) string {
	return this.Cache.Get(key)
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
