package language

import (
	"fmt"
	"runtime"
)

type Exception struct {
	code    int
	message string
	stack   string
}

func NewException(code int, message string) *Exception {
	stack := ""
	for i := 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
	}

	return &Exception{
		code:    code,
		message: message,
		stack:   stack,
	}
}

func (this *Exception) GetCode() int {
	return this.code
}

func (this *Exception) GetMessage() string {
	return this.message
}

func (this *Exception) GetStackTrace() string {
	return this.stack
}

func (this Exception) Error() string {
	return fmt.Sprintf("[Code:%d] [Message:%s] [Stack:%s]", this.code, this.message, this.stack)
}

func Throw(code int, message string) {
	panic(NewException(code, message))
}

func CatchCrash(handler func(Exception)) {
	err := recover()
	if err != nil {
		var errStr string
		exceptionErrr, isException := err.(*Exception)
		if isException {
			handler(*exceptionErrr)
		} else {
			errErr, isErr := err.(error)
			if isErr {
				errStr = errErr.Error()
			} else {
				errStr = fmt.Sprint(err)
			}
			handler(*NewException(1, errStr))
		}
	}
}

func Catch(handler func(Exception)) {
	err := recover()
	if err != nil {
		exceptionErrr, ok := err.(*Exception)
		if ok {
			handler(*exceptionErrr)
		} else {
			panic(err)
		}
	}
}
