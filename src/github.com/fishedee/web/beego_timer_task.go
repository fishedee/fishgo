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
	go func(){
		for {
			time.Sleep(duraction)
			startSingleTimerTask(handler)
		}
	}();
}
