package util

import (
	"errors"
	"reflect"
	"strconv"
	"sync"
)

const (
	//并发安全，每个数据严格地只执行一次
	MemoryFuncCacheNormal = iota
	//并发安全，每个数据可能执行多次
	MemoryFuncCacheFast
)

type memoryFuncResult struct {
	data interface{}
	err  error
}

type MemoryFunc struct {
	mutex     sync.RWMutex
	wmutex    sync.Mutex
	handler   reflect.Value
	cacheMode int
	data      map[interface{}]memoryFuncResult
}

func NewMemoryFunc(handler interface{}, cacheMode int) (*MemoryFunc, error) {
	handlerType := reflect.TypeOf(handler)
	handlerValue := reflect.ValueOf(handler)
	if handlerType.Kind() != reflect.Func {
		return nil, errors.New("memoryfunc must be a func")
	}

	if cacheMode != MemoryFuncCacheNormal &&
		cacheMode != MemoryFuncCacheFast {
		return nil, errors.New("invalid cacheMode " + strconv.Itoa(cacheMode))
	}

	result := MemoryFunc{}
	result.handler = handlerValue
	result.cacheMode = cacheMode
	result.data = map[interface{}]memoryFuncResult{}

	return &result, nil
}

func (this *MemoryFunc) handleSingleRequestByFast(request interface{}) memoryFuncResult {
	resultValue := this.handler.Call([]reflect.Value{reflect.ValueOf(request)})
	var errError error
	data := resultValue[0].Interface()
	err := resultValue[1].Interface()
	if err != nil {
		errError = err.(error)
	}
	result := memoryFuncResult{
		data: data,
		err:  errError,
	}
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.data[request] = result
	return result
}

func (this *MemoryFunc) handleSingleRequestByNormal(request interface{}) memoryFuncResult {
	this.wmutex.Lock()
	defer this.wmutex.Unlock()
	result, ok := this.hasSingleRequest(request)
	if ok {
		return result
	}
	result = this.handleSingleRequestByFast(request)
	return result
}

func (this *MemoryFunc) hasSingleRequest(request interface{}) (memoryFuncResult, bool) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	result, ok := this.data[request]
	return result, ok
}

func (this *MemoryFunc) Call(request interface{}) (interface{}, error) {
	//检查是否有数据
	result, ok := this.hasSingleRequest(request)
	if ok {
		return result.data, result.err
	}

	//获取数据
	if this.cacheMode == MemoryFuncCacheFast {
		result = this.handleSingleRequestByFast(request)
		return result.data, result.err
	} else {
		result = this.handleSingleRequestByNormal(request)
		return result.data, result.err
	}
}
