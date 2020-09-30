package language

import (
	"fmt"
	"runtime"
	"strings"
)

type Exception struct {
	code    int
	message string
	stack   []string
	cause   interface{}
	isCrash bool
}

func NewException(code int, message string, args ...interface{}) *Exception {
	return newException(2, nil, false, code, message, args...)
}

func newException(stackBegin int, cause interface{}, isCrash bool, code int, message string, args ...interface{}) *Exception {
	if len(args) != 0 {
		message = fmt.Sprintf(message, args...)
	}
	stack := []string{}
	for i := stackBegin; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		stack = append(stack, fmt.Sprintf("%s:%d", file, line))
	}

	return &Exception{
		code:    code,
		message: message,
		stack:   stack,
		cause:   cause,
	}
}

func (this *Exception) GetCode() int {
	return this.code
}

func (this *Exception) GetMessage() string {
	return this.message
}

func (this *Exception) GetCause() interface{} {
	return this.cause
}

func (this *Exception) IsCrash() bool {
	return this.isCrash
}

func (this *Exception) GetStackTrace() string {
	return strings.Join(this.stack, "\n")
}

func (this *Exception) GetStackTraceLine(i int) string {
	return this.stack[i]
}

func (this *Exception) Error() string {
	return fmt.Sprintf("[Code:%d] [Message:%s] [Stack:%s]", this.GetCode(), this.GetMessage(), this.GetStackTrace())
}

func Throw(code int, message string, args ...interface{}) {
	exception := newException(2, nil, false, code, message, args...)

	panic(exception)
}

func CatchCrash(handler func(Exception)) {
	err := recover()
	if err != nil {
		exception, isException := err.(*Exception)
		if isException {
			handler(*exception)
		} else {
			exception := newException(3, err, true, 1, fmt.Sprint(err))
			handler(*exception)
		}
	}
}

func Catch(handler func(Exception)) {
	err := recover()
	if err != nil {
		exception, isException := err.(*Exception)
		if isException {
			if exception.IsCrash() == false {
				handler(*exception)
			} else {
				panic(exception)
			}
		} else {
			exception := newException(3, err, true, 1, fmt.Sprint(err))
			panic(exception)
		}
	}
}
