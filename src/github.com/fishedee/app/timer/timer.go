package web

import (
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/language"
	"github.com/robfig/cron"
)

type TimerHandler func()

type Timer interface {
	Cron(cronspec string, handler TimerHandler) error
	MustCron(cronspec string, handler TimerHandler)

	Run() error
	Close()
}

type timerImplement struct {
	log     Log
	crontab *cron.Cron
}

func NewTimer(log Log) (Timer, error) {
	return &timerImplement{
		log:     log,
		crontab: cron.New(),
	}, nil
}

func (this *timerImplement) Cron(cronspec string, handler TimerHandler) error {
	//使用crontab标准的定时器
	// (second) (minute) (hour) (day of month) (month) (day of week, optional)
	// * * * * * *
	return this.crontab.AddFunc(cronspec, func() {
		defer CatchCrash(func(exception Exception) {
			this.log.Critical("TimeTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
		})
		handler()
	})
}

func (this *timerImplement) MustCron(cronspec string, handler TimerHandler) {
	err := this.Cron(cronspec, handler)
	if err != nil {
		panic(err)
	}
}

func (this *timerImplement) Run() error {
	this.crontab.Run()
	return nil
}

func (this *timerImplement) Close() {
	this.crontab.Stop()
}
