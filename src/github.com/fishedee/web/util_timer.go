package web

import (
	. "github.com/fishedee/language"
	"github.com/robfig/cron"
	"time"
)

type Timer interface {
	WithLog(log Log) Timer
	Cron(cronspec string, handler func()) error
	Interval(duraction time.Duration, handler func()) error
	Tick(duraction time.Duration, handler func()) error
}

type timerImplement struct {
	log Log
}

func NewTimer() (Timer, error) {
	return &timerImplement{}, nil
}

func (this *timerImplement) WithLog(log Log) Timer {
	result := *this
	result.log = log
	return &result
}

func (this *timerImplement) startSingleTask(handler func()) {
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
	crontab.Start()
	return nil
}

func (this *timerImplement) Interval(duraction time.Duration, handler func()) error {
	//带有延后属性的定时器
	go func() {
		timeChan := time.After(duraction)
		for {
			<-timeChan
			this.startSingleTask(handler)
			timeChan = time.After(duraction)
		}
	}()
	return nil
}

func (this *timerImplement) Tick(duraction time.Duration, handler func()) error {
	//带有延后属性的定时器
	go func() {
		tickChan := time.Tick(duraction)
		for {
			<-tickChan
			this.startSingleTask(handler)
		}
	}()
	return nil
}
