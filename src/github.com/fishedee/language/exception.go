package language

import (
	"fmt"
	"runtime"
)

type Exception struct{
	code int
	message string
	stack string
}

func NewException(code int,message string)(error){
	var stack string
	for i := 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
	}

	return &Exception{
		code:code,
		message:message,
		stack:stack,
	}
}

func (this *Exception)Error()(string){
	return this.message
}

func (this *Exception)GetCode()(int){
	return this.code
}

func (this *Exception)GetMessage()(string){
	return this.message
}

func (this *Exception)GetStackTrace()(string){
	return this.stack
}

func Throw(code int,message string){
	panic(NewException(code,message))
}

func Catch(handler func(Exception)){
	err := recover();
	if err != nil{
		exceptionErrr,ok := err.(*Exception)
		if ok{
			handler(*exceptionErrr)
		}else{
			panic(err)
		}
	}
}