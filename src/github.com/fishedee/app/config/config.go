package config

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/config"
	. "github.com/fishedee/language"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Config interface {
	String(key string) (string, error)
	MustString(key string) string

	StringList(key string) ([]string, error)
	MustStringList(key string) []string

	Float(key string) (float64, error)
	MustFloat(key string) float64

	Int(key string) (int, error)
	MustInt(key string) int

	Bool(key string) (bool, error)
	MustBool(key string) bool

	Duration(key string) (time.Duration, error)
	MustDuration(key string) time.Duration

	Bind(prefix string, data interface{}) error
	MustBind(prefix string, data interface{})
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
	} else if configRunMode := configer.String("runmode"); configRunMode != "" {
		runMode = configRunMode
	} else {
		runMode = "dev"
	}

	return &configImplement{
		runMode:  runMode,
		configer: configer,
	}, nil
}

func (this *configImplement) String(key string) (string, error) {
	if strings.ToLower(key) == "runmode" {
		return this.runMode, nil
	}
	if v := this.configer.String(this.runMode + "::" + key); v != "" {
		return v, nil
	}
	return this.configer.String(key), nil
}

func (this *configImplement) MustString(key string) string {
	result, err := this.String(key)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *configImplement) StringList(key string) ([]string, error) {
	v, err := this.String(key)
	if err != nil {
		return nil, err
	}
	return Explode(v, ","), nil
}

func (this *configImplement) MustStringList(key string) []string {
	result, err := this.StringList(key)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *configImplement) Float(key string) (float64, error) {
	v, err := this.String(key)
	if err != nil {
		return 0, err
	}
	if v == "" {
		return 0, nil
	}
	return strconv.ParseFloat(v, 64)
}

func (this *configImplement) MustFloat(key string) float64 {
	result, err := this.Float(key)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *configImplement) Int(key string) (int, error) {
	v, err := this.String(key)
	if err != nil {
		return 0, err
	}
	if v == "" {
		return 0, nil
	}
	return strconv.Atoi(v)
}

func (this *configImplement) MustInt(key string) int {
	result, err := this.Int(key)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *configImplement) Bool(key string) (bool, error) {
	v, err := this.String(key)
	if err != nil {
		return false, err
	}
	if v == "" {
		return false, nil
	}
	return strconv.ParseBool(v)
}

func (this *configImplement) MustBool(key string) bool {
	result, err := this.Bool(key)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *configImplement) Duration(key string) (time.Duration, error) {
	v, err := this.String(key)
	if err != nil {
		return 0, err
	}
	if v == "" {
		return 0, nil
	}
	return time.ParseDuration(v)
}

func (this *configImplement) MustDuration(key string) time.Duration {
	result, err := this.Duration(key)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *configImplement) Bind(prefix string, data interface{}) error {
	structInfo := ArrayToMap(reflect.ValueOf(data).Elem().Interface(), "config").(map[string]interface{})
	for key, value := range structInfo {
		var err error
		prefixKey := prefix + key
		if _, isOk := value.(string); isOk {
			structInfo[key], err = this.String(prefixKey)
		} else if _, isOk := value.(float64); isOk {
			structInfo[key], err = this.Float(prefixKey)
		} else if _, isOk := value.(int); isOk {
			structInfo[key], err = this.Int(prefixKey)
		} else if _, isOk := value.(bool); isOk {
			structInfo[key], err = this.Bool(prefixKey)
		} else if _, isOk := value.(time.Duration); isOk {
			structInfo[key], err = this.Duration(prefixKey)
		} else if _, isOk := value.([]interface{}); isOk {
			structInfo[key], err = this.StringList(prefixKey)
		} else {
			err = errors.New("invalid type of structInfo: " + prefixKey)
		}
		if err != nil {
			return fmt.Errorf("[key:%v][error:%v]", key, err)
		}
	}
	err := MapToArray(structInfo, data, "config")
	if err != nil {
		return err
	}
	return nil
}

func (this *configImplement) MustBind(prefix string, data interface{}) {
	err := this.Bind(prefix, data)
	if err != nil {
		panic(err)
	}
}
