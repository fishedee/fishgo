package util

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
	"strconv"
	"strings"
	"time"
)

type CacheManagerConfig struct {
	Driver     string
	SavePath   string
	SavePrefix string
	GcInterval int
}

type CacheManager struct {
	store      cache.Cache
	saveprefix string
	Log        *LogManager
}

var newCacheManagerMemory *MemoryFunc
var newCacheManagerFromConfigMemory *MemoryFunc

func init() {
	var err error
	newCacheManagerMemory, err = NewMemoryFunc(newCacheManager, MemoryFuncCacheNormal)
	if err != nil {
		panic(err)
	}
	newCacheManagerFromConfigMemory, err = NewMemoryFunc(newCacheManagerFromConfig, MemoryFuncCacheNormal)
	if err != nil {
		panic(err)
	}
}

func newCacheManager(config CacheManagerConfig) (*CacheManager, error) {
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
		return &CacheManager{
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
		return &CacheManager{
			store:      cacheInner,
			saveprefix: config.SavePrefix,
		}, nil
	} else {
		return nil, errors.New("invalid cache config " + config.Driver)
	}
}

func NewCacheManager(config CacheManagerConfig) (*CacheManager, error) {
	result, err := newCacheManagerMemory.Call(config)
	if err != nil {
		return nil, err
	}
	return result.(*CacheManager), err
}

func newCacheManagerFromConfig(configName string) (*CacheManager, error) {
	driver := beego.AppConfig.String(configName + "driver")
	savepath := beego.AppConfig.String(configName + "savepath")
	saveprefix := beego.AppConfig.String(configName + "saveprefix")
	gcintervalStr := beego.AppConfig.String(configName + "gcinterval")
	gcinterval, _ := strconv.Atoi(gcintervalStr)

	cacheConfig := CacheManagerConfig{}
	cacheConfig.Driver = driver
	cacheConfig.SavePath = savepath
	cacheConfig.SavePrefix = saveprefix
	cacheConfig.GcInterval = gcinterval
	return NewCacheManager(cacheConfig)
}

func NewCacheManagerFromConfig(configName string) (*CacheManager, error) {
	result, err := newCacheManagerFromConfigMemory.Call(configName)
	if err != nil {
		return nil, err
	}
	return result.(*CacheManager), err
}

func NewCacheManagerWithLog(log *LogManager, cache *CacheManager) *CacheManager {
	if cache == nil {
		return nil
	} else {
		newCache := *cache
		newCache.Log = log
		return &newCache
	}
}

func (this *CacheManager) Get(key string) (string, bool) {
	result := this.store.Get(this.saveprefix + key)
	if result == nil {
		return "", false
	}
	return string(result.([]byte)), true
}

func (this *CacheManager) Set(key string, value string, timeout time.Duration) {
	defer CatchCrash(func(exception Exception) {
		this.Log.Critical("Cache Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	err := this.store.Put(this.saveprefix+key, []byte(value), timeout)
	if err != nil {
		panic(err)
	}
}

func (this *CacheManager) Del(key string) {
	defer CatchCrash(func(exception Exception) {
		this.Log.Critical("Cache Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	err := this.store.Delete(this.saveprefix + key)
	if err != nil && strings.Index(err.Error(), "not exist") == -1 {
		panic(err)
	}
}
