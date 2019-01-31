package testdata

import (
	. "github.com/fishedee/language"
	"time"
)

type User struct {
	UserId     int
	Age        int
	Name       string
	CreateTime time.Time
}

type Sex struct {
	IsMale bool
}

func logic() {
	QueryColumn([]User{}, "UserId")
	QuerySelect([]User{}, func(d User) Sex {
		return Sex{}
	})
	QueryWhere([]User{}, func(c User) bool {
		return true
	})
	QuerySort([]User{}, "UserId desc,Name asc,CreateTime asc")
	QuerySort([]User{}, "UserId asc")
	QueryColumnMap([]User{}, "UserId")
}
