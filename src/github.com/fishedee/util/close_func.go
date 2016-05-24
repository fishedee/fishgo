package util

import (
	"sync"
)

type CloseFunc struct {
	mutex        sync.Mutex
	waitGroup    sync.WaitGroup
	choseHandler []func()
}

func NewCloseFunc() *CloseFunc {
	return &CloseFunc{
		choseHandler: []func(){},
	}
}

func (this *CloseFunc) AddCloseHandler(handler func()) {
	this.mutex.Lock()
	this.choseHandler = append(this.choseHandler, handler)
	this.mutex.Unlock()
}

func (this *CloseFunc) IncrCloseCounter() {
	this.waitGroup.Add(1)
}

func (this *CloseFunc) DecrCloseCounter() {
	this.waitGroup.Done()
}

func (this *CloseFunc) Close() {
	//复制一份
	var result = []func(){}
	this.mutex.Lock()
	for _, singleHandler := range this.choseHandler {
		result = append(result, singleHandler)
	}
	this.mutex.Unlock()

	//执行closeHandler
	for _, singleHandler := range result {
		this.waitGroup.Add(1)
		go func() {
			defer this.waitGroup.Done()
			singleHandler()
		}()
	}

	this.waitGroup.Wait()
}
