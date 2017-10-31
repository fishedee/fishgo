package util

import (
	"sync"
)

type LocalCache struct {
	data   map[string]interface{}
	trySet map[string]bool
	mutex  sync.RWMutex
}

func NewLocalCache() *LocalCache {
	cache := LocalCache{}
	cache.data = map[string]interface{}{}
	cache.trySet = map[string]bool{}
	return &cache
}

func (this *LocalCache) Set(key string, handler func(key string) interface{}) {
	this.mutex.Lock()
	result := this.data[key]
	if result != nil {
		this.mutex.Unlock()
		return
	}
	hasTry := this.trySet[key]
	if hasTry {
		this.mutex.Unlock()
		return
	}
	this.trySet[key] = true
	this.mutex.Unlock()

	data := handler(key)

	this.mutex.Lock()
	this.data[key] = data
	this.trySet[key] = false
	this.mutex.Unlock()
}

func (this *LocalCache) Get(key string) interface{} {
	this.mutex.RLock()
	result := this.data[key]
	this.mutex.RUnlock()

	return result
}

func (this *LocalCache) Size() int {
	this.mutex.RLock()
	size := len(this.data)
	this.mutex.RUnlock()
	return size
}

func (this *LocalCache) BatchGet(key []string) map[string]interface{} {
	result := map[string]interface{}{}

	this.mutex.RLock()
	for _, singleKey := range key {
		singleData := this.data[singleKey]
		result[singleKey] = singleData
	}
	this.mutex.RUnlock()

	return result
}
