package web

import (
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
	"github.com/robfig/cron"
	"time"
)

type Timer interface {
	WithLog(log Log) Timer
	Cron(cronspec string, handler func()) error
	Interval(duraction time.Duration, handler func()) error
	Tick(duraction time.Duration, handler func()) error
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

func (this *timerImplement) Cron(cronspec string, handler func()) error {
	//使用crontab标准的定时器
	// (second) (minute) (hour) (day of month) (month) (day of week, optional)
	// * * * * * *
	var err error
	crontab := cron.New()
	err = crontab.AddFunc(cronspec, func() {
		this.startSingleTask(handler)
	})
	if err != nil {
		return err
	}
	this.closeFunc.AddCloseHandler(func() {
		crontab.Stop()
	})
	crontab.Start()
	return nil
}

func (this *timerImplement) Interval(duraction time.Duration, handler func()) error {
	//带有延后属性的定时器
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
	return nil
}

func (this *timerImplement) Tick(duraction time.Duration, handler func()) error {
	//带有延后属性的定时器
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
	return nil
}

func (this *timerImplement) Close() {
	this.closeFunc.Close()
}
