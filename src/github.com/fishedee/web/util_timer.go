package util

import (
	. "github.com/fishedee/language"
	"github.com/robfig/cron"
	"time"
)

type TimerManager struct {
	Log *LogManager
}

func NewTimerManager() (*TimerManager, error) {
	return &TimerManager{}, nil
}

func NewTimerManagerWithLog(log *LogManager, manager *TimerManager) *TimerManager {
	if manager == nil {
		return nil
	} else {
		return &TimerManager{
			Log: log,
		}
	}
}

func (this *TimerManager) startSingleTask(handler func()) {
	defer CatchCrash(func(exception Exception) {
		this.Log.Critical("TimerTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	defer Catch(func(exception Exception) {
		this.Log.Error("TimerTask Error Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	handler()
}

func (this *TimerManager) Cron(cronspec string, handler func()) error {
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

func (this *TimerManager) Interval(duraction time.Duration, handler func()) error {
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

func (this *TimerManager) Tick(duraction time.Duration, handler func()) error {
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
