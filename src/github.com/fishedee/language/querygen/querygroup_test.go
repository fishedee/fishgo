package main

import (
	. "github.com/fishedee/assert"
	. "github.com/fishedee/language"
	. "github.com/fishedee/language/querygen/testdata"
	"math/rand"
	"testing"
	"time"
)

func TestQueryGroup(t *testing.T) {
	data := []User{
		User{UserId: 3, Name: "a"},
		User{UserId: 3, Name: "c"},
		User{UserId: 23, Name: "d"},
		User{UserId: 23, Name: "c", CreateTime: time.Unix(29, 0)},
		User{UserId: 23, Name: "c", CreateTime: time.Unix(1, 0)},
		User{UserId: 23, Name: "c", CreateTime: time.Unix(33, 0)},
		User{UserId: 23, Name: "a"},
		User{UserId: 1},
		User{UserId: 1},
	}
	AssertEqual(t, QueryGroup(data, "UserId", func(users []User) Department {
		return Department{
			Employees: users,
		}
	}), []Department{
		Department{Employees: []User{
			User{UserId: 3, Name: "a"},
			User{UserId: 3, Name: "c"},
		}},
		Department{Employees: []User{
			User{UserId: 23, Name: "d"},
			User{UserId: 23, Name: "c", CreateTime: time.Unix(29, 0)},
			User{UserId: 23, Name: "c", CreateTime: time.Unix(1, 0)},
			User{UserId: 23, Name: "c", CreateTime: time.Unix(33, 0)},
			User{UserId: 23, Name: "a"},
		}},
		Department{Employees: []User{
			User{UserId: 1},
			User{UserId: 1},
		}},
	})
	AssertEqual(t, QueryGroup([]int{1, 3, 4, 4, 3, 3}, ".", func(ids []int) Department {
		users := QuerySelect(ids, func(id int) User {
			return User{UserId: id}
		}).([]User)
		return Department{Employees: users}
	}), []Department{
		Department{Employees: []User{
			User{UserId: 1},
		}},
		Department{Employees: []User{
			User{UserId: 3},
			User{UserId: 3},
			User{UserId: 3},
		}},
		Department{Employees: []User{
			User{UserId: 4},
			User{UserId: 4},
		}},
	})
}

func initQueryGroupData() []User {
	data := make([]User, 1000, 1000)
	for i := range data {
		data[i].UserId = rand.Int()
		data[i].Age = rand.Int()
	}
	return data
}

func BenchmarkQueryGroupHand(b *testing.B) {
	data := initQueryGroupData()

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		findMap := make(map[int][]User, len(data))
		result := make([]Department, 0, len(data))
		for _, single := range data {
			users, isExist := findMap[single.UserId]
			if isExist == false {
				users = []User{}
			}
			users = append(users, single)
			findMap[single.UserId] = users
		}
		for _, single := range data {
			users, isExist := findMap[single.UserId]
			if isExist {
				continue
			}
			delete(findMap, single.UserId)
			result = append(result, Department{
				Employees: users,
			})
		}
	}
}

func BenchmarkQueryGroupMacro(b *testing.B) {
	data := initQueryGroupData()

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		QueryGroup(data, "UserId", func(users []User) Department {
			return Department{
				Employees: users,
			}
		})
	}
}

func BenchmarkQueryGroupReflect(b *testing.B) {
	data := initQueryGroupData()

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		QueryGroup(data, "Age", func(users []User) Department {
			return Department{
				Employees: users,
			}
		})
	}
}
