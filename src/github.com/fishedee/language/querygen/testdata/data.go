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
	IsMale  bool
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
	QueryColumn([]User{}, ".")
	QueryColumn([]int{}, ".")
	QuerySelect([]User{}, func(d User) Sex {
		return Sex{}
	})
	QueryWhere([]int{}, func(c int) bool {
		return c%2 == 0
	})
	QueryWhere([]User{}, func(c User) bool {
		return true
	})
	QuerySort([]User{}, "UserId desc,Name asc,CreateTime asc")
	QuerySort([]User{}, "UserId asc")
	QuerySort([]Admin{}, "IsMale asc")
	QuerySort([]int{}, ". desc")
	QueryColumnMap([]User{}, "UserId")
	QueryColumnMap([]User{}, "[]UserId")
	QueryColumnMap([]int{}, ".")
	QueryGroup([]User{}, "UserId", func(user []User) Department {
		return Department{}
	})
	QueryGroup([]User{}, "CreateTime", func(user []User) Department {
		return Department{}
	})
	QueryGroup([]User{}, "CreateTime", func(user []User) []Department {
		return []Department{}
	})
	QueryGroup([]int{}, ".", func(ids []int) Department {
		users := QuerySelect(ids, func(id int) User {
			return User{UserId: id}
		}).([]User)
		return Department{Employees: users}
	})
	QueryLeftJoin([]Admin{}, []User{}, "AdminId = UserId", func(left Admin, right User) AdminUser {
		return AdminUser{}
	})
	QueryRightJoin([]User{}, []int{}, "UserId = .", func(left User, right int) User {
		return User{}
	})
	QueryJoin([]Admin{}, []User{}, "inner", "AdminId = UserId", func(left Admin, right User) AdminUser {
		return AdminUser{}
	})
	QueryCombine([]Admin{}, []User{}, func(left Admin, right User) AdminUser {
		return AdminUser{}
	})
	QueryCombine([]int{}, []User{}, func(left int, right User) User {
		return User{}
	})
}
