package web

import (
	"time"
	. "github.com/fishedee/language"
)

func startSingleTimerTask(handler func()){
	defer CatchCrash(func(exception Exception){
		Log.Critical("TimerTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]",exception.GetCode(),exception.GetMessage(),exception.GetStackTrace())
	})
	defer Catch(func(exception Exception){
		Log.Error("TimerTask Error Code:[%d] Message:[%s]\nStackTrace:[%s]",exception.GetCode(),exception.GetMessage(),exception.GetStackTrace())
	})
	handler()
}

func StartTimerTask(duraction time.Duration,handler func() ){
	//带有延后属性的定时器
	go func(){
		timeChan := time.After(duraction)
		for {
			<- timeChan
			startSingleTimerTask(handler)
			timeChan = time.After(duraction)
		}
	}();
}

func StartTickTask(duraction time.Duration,handler func() ){
	//没有延后属性的定时器
	go func(){
		tickChan := time.Tick(duraction)
		for {
			<- tickChan
			startSingleTimerTask(handler)
		}
	}();
}