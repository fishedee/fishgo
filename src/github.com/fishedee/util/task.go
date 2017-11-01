package util

import (
	"reflect"
	"sync"
)

type taskData struct {
	isDone bool
	data   []interface{}
}

type Task struct {
	threadCount     int
	bufferCount     int
	taskChannel     chan taskData
	doneChannel     []chan bool
	handler         reflect.Value
	progressHandler func(int, int)
	isRunning       bool
	isAutoStop      bool
	allCount        int
	finishCount     int
	progressMutex   *sync.Mutex
	countMutex      *sync.Mutex
	taskLock        sync.Locker
	taskCond        *sync.Cond
	taskList        []taskData
}

func NewTask() *Task {
	task := Task{}
	return &task
}

func (this *Task) GetIsRunning() bool {
	return this.isRunning
}

func (this *Task) GetAllCount() int {
	return this.allCount
}

func (this *Task) GetFinishCount() int {
	return this.finishCount
}

func (this *Task) SetIsAutoStop(isAutoStop bool) {
	if this.isRunning {
		return
	}
	this.isAutoStop = isAutoStop
}

func (this *Task) SetBufferCount(bufferCount int) {
	if this.isRunning {
		return
	}
	this.bufferCount = bufferCount
}

func (this *Task) SetThreadCount(threadCount int) {
	if this.isRunning {
		return
	}
	this.threadCount = threadCount
}

func (this *Task) SetHandler(handler interface{}) {
	if this.isRunning {
		return
	}
	this.handler = reflect.ValueOf(handler)
}

func (this *Task) SetProgressHandler(progressHandler func(int, int)) {
	if this.isRunning {
		return
	}
	this.progressHandler = progressHandler
}

func (this *Task) AddTask(data ...interface{}) {
	if !this.isRunning {
		return
	}
	this.countMutex.Lock()
	this.allCount++
	this.countMutex.Unlock()

	this.taskCond.L.Lock()
	this.taskList = append(this.taskList, taskData{
		isDone: false,
		data:   data,
	})
	this.taskCond.Signal()
	this.taskCond.L.Unlock()
}

func (this *Task) Start() {
	if this.isRunning {
		return
	}
	this.allCount = 0
	this.finishCount = 0
	this.progressMutex = &sync.Mutex{}
	this.countMutex = &sync.Mutex{}
	this.taskChannel = make(chan taskData, this.bufferCount)
	this.doneChannel = make([]chan bool, this.threadCount)
	this.taskLock = &sync.Mutex{}
	this.taskCond = sync.NewCond(this.taskLock)
	this.isRunning = true
	go func() {
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
				this.taskChannel <- singleTask
			}
			if isStop {
				for i := 0; i != this.threadCount; i++ {
					this.taskChannel <- taskData{
						isDone: true,
						data:   nil,
					}
				}
				break
			}
		}
	}()

	for i := 0; i != this.threadCount; i++ {
		this.doneChannel[i] = make(chan bool)
		curDoneChannel := this.doneChannel[i]
		go func() {
			for {
				task := <-this.taskChannel
				if task.isDone {
					break
				}
				taskData := []reflect.Value{}
				for _, singleData := range task.data {
					taskData = append(taskData, reflect.ValueOf(singleData))
				}

				func() {
					this.handler.Call(taskData)
				}()

				this.countMutex.Lock()
				this.finishCount++
				this.countMutex.Unlock()
				if this.isAutoStop && this.finishCount == this.allCount {
					this.Stop()
				}
				if this.progressHandler != nil {
					this.progressMutex.Lock()
					this.progressHandler(this.finishCount, this.allCount)
					this.progressMutex.Unlock()
				}
			}
			curDoneChannel <- true
		}()
	}
}

func (this *Task) Wait() {
	for i := 0; i != this.threadCount; i++ {
		<-this.doneChannel[i]
	}
}

func (this *Task) Stop() {
	if !this.isRunning {
		return
	}
	this.isRunning = false
	this.taskCond.L.Lock()
	this.taskList = append(this.taskList, taskData{
		isDone: true,
		data:   nil,
	})
	this.taskCond.Signal()
	this.taskCond.L.Unlock()
}
