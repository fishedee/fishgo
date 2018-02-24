package web

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	"reflect"
	"strings"
	"sync"
	"time"
)

type Cache interface {
	Get(key string) (string, error)
	MustGet(key string) string

	Set(key string, value string, timeout time.Duration) error
	MustSet(key string, value string, timeout time.Duration)

	Delete(key string) error
	MustDelete(key string)

	Memoize(key string, value interface{}, timeout time.Duration) (interface{}, error)
	MustMemoize(key string, valuer interface{}, timeout time.Duration) interface{}
}

type CacheConfig struct {
	Driver     string `config::"dirver"`
	SavePath   string `config::"savepath"`
	SavePrefix string `config::"saveprefix"`
	GcInterval int    `config::"gcinterval"`
}

type cacheHandler struct {
	decode func([]byte) (interface{}, error)
	encode func(interface{}) ([]byte, error)
}

type cacheImplement struct {
	store      cache.Cache
	saveprefix string
	typeInfo   sync.Map
}

func NewCache(config CacheConfig) (Cache, error) {
	if config.Driver == "" {
		return nil, nil
	} else if config.Driver == "memory" {
		var data struct {
			Interval int `json:"interval,omitempty"`
		}
		data.Interval = config.GcInterval
		configString, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		cacheInner, err := cache.NewCache("memory", string(configString))
		if err != nil {
			return nil, err
		}
		return &cacheImplement{
			store:      cacheInner,
			saveprefix: config.SavePrefix,
		}, nil
	} else if config.Driver == "redis" {
		var data struct {
			Key      string `json:"key"`
			Conn     string `json:"conn"`
			Password string `json:"password,omitempty"`
		}
		if config.SavePrefix == "" {
			return nil, errors.New("invalid config.SavePrefix is empty")
		}
		data.Key = config.SavePrefix
		configArray := Explode(config.SavePath, ",")
		if len(configArray) == 0 {
			return nil, errors.New("invalid config.SavePath " + config.SavePath)
		}
		data.Conn = configArray[0]
		if len(configArray) >= 3 {
			data.Password = configArray[2]
		}
		configString, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		cacheInner, err := cache.NewCache("redis", string(configString))
		if err != nil {
			return nil, err
		}
		return &cacheImplement{
			store:      cacheInner,
			saveprefix: config.SavePrefix,
		}, nil
	} else {
		return nil, errors.New("invalid cache config " + config.Driver)
	}
}

func (this *cacheImplement) getInner(key string) ([]byte, error) {
	result := this.store.Get(this.saveprefix + key)
	if result == nil {
		return nil, nil
	}
	return result.([]byte), nil
}

func (this *cacheImplement) setInner(key string, value []byte, timeout time.Duration) error {
	err := this.store.Put(this.saveprefix+key, value, timeout)
	if err != nil {
		return err
	}
	return nil
}

func (this *cacheImplement) Get(key string) (string, error) {
	result, err := this.getInner(key)
	return string(result), err
}

func (this *cacheImplement) MustGet(key string) string {
	result, err := this.Get(key)
	if err != nil {
		panic(err)
	}
	return result
}

func (this *cacheImplement) Set(key string, value string, timeout time.Duration) error {
	err := this.setInner(key, []byte(value), timeout)
	return err
}

func (this *cacheImplement) MustSet(key string, value string, timeout time.Duration) {
	err := this.Set(key, value, timeout)
	if err != nil {
		panic(err)
	}
}

func (this *cacheImplement) Delete(key string) error {
	err := this.store.Delete(this.saveprefix + key)
	if err != nil && strings.Index(err.Error(), "not exist") == -1 {
		return err
	}
	return nil
}

func (this *cacheImplement) MustDelete(key string) {
	err := this.Delete(key)
	if err != nil {
		panic(err)
	}
}

func (this *cacheImplement) getHandler(typeInfo reflect.Type) (cacheHandler, error) {
	handler, isExist := this.typeInfo.Load(typeInfo)
	if isExist {
		return handler.(cacheHandler), nil
	}
	if typeInfo.Kind() != reflect.Func {
		return cacheHandler{}, errors.New("invalid Memoize value ,must be func")
	}
	if typeInfo.NumIn() != 0 {
		return cacheHandler{}, errors.New("invalid Memoize value ,must be zero argument in")
	}
	if typeInfo.NumOut() != 1 {
		return cacheHandler{}, errors.New("invalid Memoize value ,must be only one return value out")
	}
	numOut := typeInfo.Out(0)
	result := cacheHandler{
		decode: func(data []byte) (interface{}, error) {
			result := reflect.New(numOut)
			err := DecodeJson(data, result.Interface())
			if err != nil {
				return nil, err
			}
			return result.Elem().Interface(), nil
		},
		encode: func(data interface{}) ([]byte, error) {
			return EncodeJson(data)
		},
	}
	this.typeInfo.Store(typeInfo, result)
	return result, nil
}

func (this *cacheImplement) Memoize(key string, value interface{}, timeout time.Duration) (interface{}, error) {
	handler, err := this.getHandler(reflect.TypeOf(value))
	if err != nil {
		return nil, err
	}
	existValue, err := this.getInner(key)
	if err != nil {
		return nil, err
	}
	if existValue != nil {
		//已存在数据
		return handler.decode(existValue)
	} else {
		//不存在数据
		valueCall := reflect.ValueOf(value)
		callResult := valueCall.Call(nil)
		getResult := callResult[0].Interface()
		data, err := handler.encode(getResult)
		if err != nil {
			return nil, err
		}
		this.setInner(key, data, timeout)
		return getResult, nil
	}
}

func (this *cacheImplement) MustMemoize(key string, value interface{}, timeout time.Duration) interface{} {
	result, err := this.Memoize(key, value, timeout)
	if err != nil {
		panic(err)
	}
	return result
}
