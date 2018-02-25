package workgroup

import (
	"errors"
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/assert"
	"testing"
	"time"
)

type taskStub struct {
	runHandler   func() error
	closeHandler func()
}

func (this *taskStub) Run() error {
	return this.runHandler()
}

func (this *taskStub) Close() {
	this.closeHandler()
}

func TestWorkGroupError(t *testing.T) {
	log, _ := NewLog(LogConfig{
		Driver: "console",
	})
	workgroup, _ := NewWorkGroup(log, WorkGroupConfig{
		CloseTimeout: time.Second * 3,
	})

	workgroup.Add(&taskStub{
		runHandler: func() error {
			time.Sleep(time.Second)
			return errors.New("error1")
		},
		closeHandler: func() {
		},
	})
	workgroup.Add(&taskStub{
		runHandler: func() error {
			time.Sleep(time.Second * 2)
			return errors.New("error2")
		},
		closeHandler: func() {
		},
	})
	begin := time.Now().Unix()
	err := workgroup.Run()
	end := time.Now().Unix()
	AssertEqual(t, end-begin >= 2, true)
	AssertEqual(t, err != nil, true)
	AssertEqual(t, err.Error(), "error1")
}

func TestWorkGroupClose(t *testing.T) {
	log, _ := NewLog(LogConfig{
		Driver: "console",
	})
	workgroup, _ := NewWorkGroup(log, WorkGroupConfig{})

	exitChan1 := make(chan bool)
	workgroup.Add(&taskStub{
		runHandler: func() error {
			<-exitChan1
			return nil
		},
		closeHandler: func() {
			exitChan1 <- true
		},
	})

	exitChan2 := make(chan bool)
	workgroup.Add(&taskStub{
		runHandler: func() error {
			<-exitChan2
			return nil
		},
		closeHandler: func() {
			exitChan2 <- true
		},
	})

	go func() {
		time.Sleep(time.Second * 2)
		workgroup.Close()
	}()
	begin := time.Now().Unix()
	err := workgroup.Run()
	end := time.Now().Unix()
	AssertEqual(t, err == nil, true)
	AssertEqual(t, end-begin >= 2, true)
}

func TestWorkGroupTimeout(t *testing.T) {
	log, _ := NewLog(LogConfig{
		Driver: "console",
	})
	workgroup, _ := NewWorkGroup(log, WorkGroupConfig{
		CloseTimeout: time.Second * 2,
	})

	taskChan := make(chan bool)
	workgroup.Add(&taskStub{
		runHandler: func() error {
			<-taskChan
			return nil
		},
		closeHandler: func() {
		},
	})
	workgroup.Add(&taskStub{
		runHandler: func() error {
			<-taskChan
			return nil
		},
		closeHandler: func() {
		},
	})
	go func() {
		time.Sleep(time.Second)
		workgroup.Close()
	}()
	begin := time.Now().Unix()
	err := workgroup.Run()
	end := time.Now().Unix()
	AssertEqual(t, end-begin >= 3, true)
	AssertEqual(t, err == nil, true)
}

func TestWorkGroupCloseTwo(t *testing.T) {
	log, _ := NewLog(LogConfig{
		Driver: "console",
	})
	workgroup, _ := NewWorkGroup(log, WorkGroupConfig{
		CloseTimeout: time.Second,
	})

	taskChan := make(chan bool)
	workgroup.Add(&taskStub{
		runHandler: func() error {
			<-taskChan
			return nil
		},
		closeHandler: func() {
		},
	})
	exitChan := make(chan bool)
	go func() {
		time.Sleep(time.Second)
		workgroup.Close()
		workgroup.Close()
		exitChan <- true
	}()
	begin := time.Now().Unix()
	err := workgroup.Run()
	end := time.Now().Unix()
	AssertEqual(t, end-begin >= 1, true)
	AssertEqual(t, <-exitChan, true)
	AssertEqual(t, err == nil, true)
}
