package test

import (
	"time"
)

type ConfigAoModel struct {
	BaseModel
	data       map[string]string
	dataStruct map[string]ConfigData
}

type ConfigData struct {
	Data       string
	CreateTime time.Time `xorm:"created"`
	ModifyTime time.Time `xorm:"updated"`
}

func (this *ConfigAoModel) Set(key string, value string) {
	if this.data == nil {
		this.data = map[string]string{}
	}
	this.data[key] = value
}

func (this *ConfigAoModel) Get(key string) string {
	return this.data[key]
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
