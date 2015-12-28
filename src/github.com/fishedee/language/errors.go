package language

import (
	"github.com/go-errors/errors"
)

type Error struct {
	innerError *errors.Error
	code       int
	message    string
}

func NewError(code int, message string) *Error {
	return &Error{
		innerError: errors.New(message),
		code:       code,
		message:    message,
	}
}

func (this *Error) GetCode() int {
	return this.code
}

func (this *Error) GetMessage() string {
	return this.message
}

func (this *Error) Error() string {
	return this.message
}

func (this *Error) ErrorStack() string {
	return this.innerError.ErrorStack()
}
