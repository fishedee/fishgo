package testdata

import (
	. "github.com/fishedee/language"
)

type User struct {
	UserId int
	Age    int
	Name   string
}

func init() {
	QueryColumn([]User{}, "UserId")
}
