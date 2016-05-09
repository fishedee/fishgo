package test

import (
	. "github.com/fishedee/web"
	"sync/atomic"
)

type counterAoModel struct {
	Model
	totalInt int32
}

func (this *counterAoModel) Incr() {
	this.totalInt++
}

func (this *counterAoModel) IncrAtomic() {
	atomic.AddInt32(&this.totalInt, 1)
}

func (this *counterAoModel) Reset() {
	this.totalInt = 0
}

func (this *counterAoModel) Get() int {
	return (int)(this.totalInt)
}
