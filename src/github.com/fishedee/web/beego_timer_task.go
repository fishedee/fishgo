package web

import (
	"time"
	. "github.com/fishedee/language"
)

func startSingleTimerTask(handler func()){
	defer Catch(func(exception Exception){
		Log.Error(exception.GetStackTrace())
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
