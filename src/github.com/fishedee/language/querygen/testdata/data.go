package testdata

import (
	. "github.com/fishedee/language"
	"time"
)

type Department struct {
	DepartmentId int
	Name         string
	Employees    []User
}

type User struct {
	UserId     int
	Age        int
	Name       string
	CreateTime time.Time
}

type Admin struct {
	AdminId int
	Level   int
}

type AdminUser struct {
	AdminId    int
	Level      int
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
	QueryGroup([]User{}, "UserId", func(user []User) Department {
		return Department{}
	})
	QueryLeftJoin([]Admin{}, []User{}, "AdminId = UserId", func(left Admin, right User) AdminUser {
		return AdminUser{}
	})
	QueryJoin([]Admin{}, []User{}, "inner", "AdminId = UserId", func(left Admin, right User) AdminUser {
		return AdminUser{}
	})
}
