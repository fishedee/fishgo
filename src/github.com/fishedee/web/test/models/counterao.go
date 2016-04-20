package test

import (
	"sync/atomic"
)

type CounterAoModel struct {
	BaseModel
	totalInt int32
}

func (this *CounterAoModel) Incr() {
	this.totalInt++
}

func (this *CounterAoModel) IncrAtomic() {
	atomic.AddInt32(&this.totalInt, 1)
}

func (this *CounterAoModel) Reset() {
	this.totalInt = 0
}

func (this *CounterAoModel) Get() int {
	return (int)(this.totalInt)
}
