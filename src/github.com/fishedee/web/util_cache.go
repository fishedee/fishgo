package web

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	. "github.com/fishedee/language"
	"strconv"
	"strings"
	"time"
)

type Cache interface {
	WithLog(log Log) Cache
	Get(key string) (string, bool)
	Set(key string, value string, timeout time.Duration)
	Del(key string)
}

type CacheConfig struct {
	Driver     string
	SavePath   string
	SavePrefix string
	GcInterval int
}

type cacheImplement struct {
	store      cache.Cache
	saveprefix string
	log        Log
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

func NewCacheFromConfig(configName string) (Cache, error) {
	driver := globalBasic.Config.GetString(configName + "driver")
	savepath := globalBasic.Config.GetString(configName + "savepath")
	saveprefix := globalBasic.Config.GetString(configName + "saveprefix")
	gcintervalStr := globalBasic.Config.GetString(configName + "gcinterval")
	gcinterval, _ := strconv.Atoi(gcintervalStr)

	cacheConfig := CacheConfig{}
	cacheConfig.Driver = driver
	cacheConfig.SavePath = savepath
	cacheConfig.SavePrefix = saveprefix
	cacheConfig.GcInterval = gcinterval
	return NewCache(cacheConfig)
}

func (this *cacheImplement) WithLog(log Log) Cache {
	if this == nil {
		return nil
	} else {
		newCache := *this
		newCache.log = log
		return &newCache
	}
}

func (this *cacheImplement) Get(key string) (string, bool) {
	result := this.store.Get(this.saveprefix + key)
	if result == nil {
		return "", false
	}
	return string(result.([]byte)), true
}

func (this *cacheImplement) Set(key string, value string, timeout time.Duration) {
	defer CatchCrash(func(exception Exception) {
		this.log.Critical("Cache Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	err := this.store.Put(this.saveprefix+key, []byte(value), timeout)
	if err != nil {
		panic(err)
	}
}

func (this *cacheImplement) Del(key string) {
	defer CatchCrash(func(exception Exception) {
		this.log.Critical("Cache Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	err := this.store.Delete(this.saveprefix + key)
	if err != nil && strings.Index(err.Error(), "not exist") == -1 {
		panic(err)
	}
}
