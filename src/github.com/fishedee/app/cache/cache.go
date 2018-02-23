package web

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	. "github.com/fishedee/language"
	"strings"
	"time"
)

type Cache interface {
	Get(key string) (string, bool)
	Set(key string, value string, timeout time.Duration) error
	Del(key string) error
}

type CacheConfig struct {
	Driver     string `config::"dirver"`
	SavePath   string `config::"savepath"`
	SavePrefix string `config::"saveprefix"`
	GcInterval int    `config::"gcinterval"`
}

type cacheImplement struct {
	store      cache.Cache
	saveprefix string
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

func (this *cacheImplement) Get(key string) (string, bool) {
	result := this.store.Get(this.saveprefix + key)
	if result == nil {
		return "", false
	}
	return string(result.([]byte)), true
}

func (this *cacheImplement) Set(key string, value string, timeout time.Duration) error {
	err := this.store.Put(this.saveprefix+key, []byte(value), timeout)
	if err != nil {
		return err
	}
	return nil
}

func (this *cacheImplement) Del(key string) error {
	err := this.store.Delete(this.saveprefix + key)
	if err != nil && strings.Index(err.Error(), "not exist") == -1 {
		return err
	}
	return nil
}
