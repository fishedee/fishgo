package util

import (
	"sync"
)

type infiniteChannelData struct {
	isDone bool
	data   interface{}
}

type InfiniteChannel struct {
	taskLock  sync.Locker
	taskCond  *sync.Cond
	taskList  []infiniteChannelData
	outChan   chan interface{}
	isRunning bool
}

func NewInfiniteChannel() *InfiniteChannel {
	taskLock := &sync.Mutex{}
	taskCond := sync.NewCond(taskLock)
	outChan := make(chan interface{})
	result := &InfiniteChannel{
		taskLock:  taskLock,
		taskCond:  taskCond,
		outChan:   outChan,
		taskList:  nil,
		isRunning: true,
	}
	go result.run()
	return result
}

func (this *InfiniteChannel) run() {
	for {
		this.taskCond.L.Lock()
		for len(this.taskList) == 0 {
			this.taskCond.Wait()
		}
		taskList := this.taskList
		this.taskList = nil
		this.taskCond.L.Unlock()

		isStop := false
		for _, singleTask := range taskList {
			if singleTask.isDone {
				isStop = true
				break
			}
			this.outChan <- singleTask.data
		}
		if isStop {
			close(this.outChan)
			break
		}
	}
	this.isRunning = false
}

func (this *InfiniteChannel) Read() <-chan interface{} {
	return this.outChan
}

func (this *InfiniteChannel) Write(data interface{}) {
	if this.isRunning == false {
		panic("invalid write in the close InfiniteChannel")
	}
	this.taskCond.L.Lock()
	this.taskList = append(this.taskList, infiniteChannelData{
		isDone: false,
		data:   data,
	})
	this.taskCond.Signal()
	this.taskCond.L.Unlock()
}

func (this *InfiniteChannel) Close() {
	this.taskCond.L.Lock()
	this.taskList = append(this.taskList, infiniteChannelData{
		isDone: true,
		data:   nil,
	})
	this.taskCond.Signal()
	this.taskCond.L.Unlock()
}
