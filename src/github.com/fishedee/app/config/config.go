package config

import (
	"errors"
	"github.com/astaxie/beego/config"
	. "github.com/fishedee/language"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Config interface {
	GetString(key string) string
	GetFloat(key string) float64
	GetInt(key string) int
	GetBool(key string) bool
	GetStruct(prefix string, data interface{})
}

type configImplement struct {
	runMode  string
	configer config.Configer
}

func checkFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	} else {
		return true
	}
}

func NewConfig(driver string, file string) (Config, error) {
	if driver != "ini" {
		return nil, errors.New("drvier is not ini!")
	}
	configer, err := config.NewConfig("ini", file)
	if err != nil {
		return nil, err
	}

	var runMode string
	if envRunMode := os.Getenv("RUNMODE"); envRunMode != "" {
		runMode = envRunMode
	} else if configRunMode := configer.String("RunMode"); configRunMode != "" {
		runMode = configRunMode
	} else {
		runMode = "dev"
	}

	return &configImplement{
		runMode:  runMode,
		configer: configer,
	}, nil
}

func (this *configImplement) GetString(key string) string {
	if strings.ToLower(key) == "runmode" {
		return this.runMode
	}
	if v := this.configer.String(this.runMode + "::" + key); v != "" {
		return v
	}
	return this.configer.String(key)
}

func (this *configImplement) GetFloat(key string) float64 {
	v := this.GetString(key)
	vF, _ := strconv.ParseFloat(v, 64)
	return vF
}

func (this *configImplement) GetInt(key string) int {
	v := this.GetString(key)
	vI, _ := strconv.ParseInt(v, 10, 64)
	return int(vI)
}

func (this *configImplement) GetBool(key string) bool {
	v := this.GetString(key)
	vB, _ := strconv.ParseBool(v)
	return bool(vB)
}

func (this *configImplement) GetDuration(key string) time.Duration {
	v := this.GetString(key)
	vD, _ := time.ParseDuration(v)
	return vD
}

func (this *configImplement) GetStruct(prefix string, data interface{}) {
	structInfo := ArrayToMap(reflect.ValueOf(data).Elem().Interface(), "config").(map[string]interface{})
	for key, value := range structInfo {
		prefixKey := prefix + key
		if _, isOk := value.(string); isOk {
			structInfo[key] = this.GetString(prefixKey)
		} else if _, isOk := value.(float64); isOk {
			structInfo[key] = this.GetFloat(prefixKey)
		} else if _, isOk := value.(int); isOk {
			structInfo[key] = this.GetInt(prefixKey)
		} else if _, isOk := value.(bool); isOk {
			structInfo[key] = this.GetBool(prefixKey)
		} else if _, isOk := value.(time.Duration); isOk {
			structInfo[key] = this.GetDuration(prefixKey)
		} else {
			panic("invalid type of structInfo: " + prefixKey)
		}
	}
	err := MapToArray(structInfo, data, "config")
	if err != nil {
		panic(err)
	}
}
