package testdata

import (
	. "github.com/fishedee/language"
)

type User struct {
	UserId int
	Age    int
	Name   string
}

type Sex struct {
	IsMale bool
}

func logic() {
	QueryColumn([]User{}, "UserId")
	QuerySelect([]User{}, func(d User) Sex {
		return Sex{}
	})
}
