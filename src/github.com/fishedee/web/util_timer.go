package web

import (
	"errors"
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
	"github.com/robfig/cron"
	"reflect"
	"time"
)

type Timer interface {
	WithLog(log Log) Timer
	Cron(cronspec string, handler interface{})
	Interval(duraction time.Duration, handler interface{})
	Tick(duraction time.Duration, handler interface{})
	Close()
}

type timerImplement struct {
	log       Log
	closeFunc *CloseFunc
}

func NewTimer() (Timer, error) {
	return &timerImplement{
		closeFunc: NewCloseFunc(),
	}, nil
}

func (this *timerImplement) WithLog(log Log) Timer {
	result := *this
	result.log = log
	return &result
}

func (this *timerImplement) startSingleTask(handler func()) {
	this.closeFunc.IncrCloseCounter()
	defer this.closeFunc.DecrCloseCounter()
	defer CatchCrash(func(exception Exception) {
		this.log.Critical("TimerTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	defer Catch(func(exception Exception) {
		this.log.Error("TimerTask Error Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	handler()
}

func (this *timerImplement) getHandler(handler interface{}) (func(), error) {
	handlerValue := reflect.ValueOf(handler)
	handlerType := reflect.TypeOf(handler)
	if handlerType.Kind() != reflect.Func {
		return nil, errors.New("invalid handler not function")
	}
	if handlerType.NumIn() != 1 {
		return nil, errors.New("invalid handler should has a parameter")
	}
	handlerFirstArgv := handlerType.In(0)
	if handlerFirstArgv.Kind() != reflect.Ptr {
		return nil, errors.New("invalid handler first parameter is not a ptr")
	}
	handlerFirstArgv = handlerFirstArgv.Elem()
	return func() {
		target := reflect.New(handlerFirstArgv)
		basic := initEmptyBasic(nil)
		injectIoc(target, basic)
		handlerValue.Call([]reflect.Value{target})
	}, nil
}

func (this *timerImplement) Cron(cronspec string, inHandler interface{}) {
	defer CatchCrash(func(exception Exception) {
		this.log.Critical("TimerTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	//使用crontab标准的定时器
	// (second) (minute) (hour) (day of month) (month) (day of week, optional)
	// * * * * * *
	var err error
	crontab := cron.New()
	handler, err := this.getHandler(inHandler)
	if err != nil {
		panic(err)
	}
	err = crontab.AddFunc(cronspec, func() {
		this.startSingleTask(handler)
	})
	if err != nil {
		panic(err)
	}
	this.closeFunc.AddCloseHandler(func() {
		crontab.Stop()
	})
	crontab.Start()
}

func (this *timerImplement) Interval(duraction time.Duration, inHandler interface{}) {
	defer CatchCrash(func(exception Exception) {
		this.log.Critical("TimerTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	//带有延后属性的定时器
	handler, err := this.getHandler(inHandler)
	if err != nil {
		panic(err)
	}

	closeEvent := make(chan bool)
	this.closeFunc.AddCloseHandler(func() {
		closeEvent <- true
	})
	go func() {
		timeChan := time.After(duraction)
		for {
			select {
			case <-timeChan:
				this.startSingleTask(handler)
				timeChan = time.After(duraction)
			case <-closeEvent:
				return
			}

		}
	}()
}

func (this *timerImplement) Tick(duraction time.Duration, inHandler interface{}) {
	defer CatchCrash(func(exception Exception) {
		this.log.Critical("TimerTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	//带有延后属性的定时器
	handler, err := this.getHandler(inHandler)
	if err != nil {
		panic(err)
	}
	closeEvent := make(chan bool)
	this.closeFunc.AddCloseHandler(func() {
		closeEvent <- true
	})
	go func() {
		tickChan := time.Tick(duraction)
		for {
			select {
			case <-tickChan:
				this.startSingleTask(handler)
			case <-closeEvent:
				return
			}
		}
	}()
}

func (this *timerImplement) Close() {
	this.closeFunc.Close()
}
