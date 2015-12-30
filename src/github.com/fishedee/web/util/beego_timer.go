package util

import (
	. "github.com/fishedee/language"
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

func (this *TimerManager) Interval(duraction time.Duration, handler TimerHandler) {
	//带有延后属性的定时器
	go func() {
		timeChan := time.After(duraction)
		for {
			<-timeChan
			this.startSingleTask(handler)
			timeChan = time.After(duraction)
		}
	}()
}

func (this *TimerManager) Tick(duraction time.Duration, handler TimerHandler) {
	//带有延后属性的定时器
	go func() {
		tickChan := time.Tick(duraction)
		for {
			<-tickChan
			this.startSingleTask(handler)
		}
	}()
}
