package modules

import (
	"time"
)

type Timer struct {
	beginTime time.Time
	endTime   time.Time
}

func NewTimer() *Timer {
	return &Timer{}
}

func (this *Timer) Start() {
	this.beginTime = time.Now()
}

func (this *Timer) Stop() {
	this.endTime = time.Now()
}

func (this *Timer) Elapsed() time.Duration {
	return time.Duration(this.endTime.UnixNano() - this.beginTime.UnixNano())
}
