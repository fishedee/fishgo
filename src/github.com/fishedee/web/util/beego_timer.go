package util

import (
	. "github.com/fishedee/language"
	"github.com/robfig/cron"
	"time"
)

type TimerManager struct {
	Log     *LogManager
	Monitor *MonitorManager
}

type TimerHandler func()

func NewTimerManager() (*TimerManager, error) {
	return &TimerManager{}, nil
}

func NewTimerManagerWithLogAndMonitor(log *LogManager, monitor *MonitorManager, manager *TimerManager) *TimerManager {
	if manager == nil {
		return nil
	} else {
		return &TimerManager{
			Log:     log,
			Monitor: monitor,
		}
	}
}

func (this *TimerManager) startSingleTask(handler TimerHandler) {
	defer CatchCrash(func(exception Exception) {
		this.Log.Critical("TimerTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
		if this.Monitor != nil {
			this.Monitor.AscCriticalCount()
		}
	})
	defer Catch(func(exception Exception) {
		this.Log.Error("TimerTask Error Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
		if this.Monitor != nil {
			this.Monitor.AscErrorCount()
		}
	})
	handler()
}

func (this *TimerManager) Cron(cronspec string, handler TimerHandler) error {
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

func (this *TimerManager) Interval(duraction time.Duration, handler TimerHandler) error {
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

func (this *TimerManager) Tick(duraction time.Duration, handler TimerHandler) error {
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
