package main

import (
	. "github.com/fishedee/assert"
	. "github.com/fishedee/language"
	. "github.com/fishedee/language/querygen/testdata"
	"math/rand"
	"sort"
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
			User{UserId: 1},
			User{UserId: 1},
		}},
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
	})
	AssertEqual(t, QueryGroup([]int{2, 3, 4, 6, 2, 2, 9, 6}, ". desc", func(ids []int) Department {
		users := QuerySelect(ids, func(id int) User {
			return User{UserId: id}
		}).([]User)
		return Department{Employees: users}
	}), []Department{
		Department{Employees: []User{
			User{UserId: 9},
		}},
		Department{Employees: []User{
			User{UserId: 6},
			User{UserId: 6},
		}},
		Department{Employees: []User{
			User{UserId: 4},
		}},
		Department{Employees: []User{
			User{UserId: 3},
		}},
		Department{Employees: []User{
			User{UserId: 2},
			User{UserId: 2},
			User{UserId: 2},
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
		newData := make([]User, len(data), len(data))
		newData2 := make([]Department, 0, len(data))
		copy(newData, data)
		sort.SliceStable(newData, func(i int, j int) bool {
			return newData[i].UserId < newData[j].UserId
		})
		for i := 0; i != len(newData); i++ {
			j := i
			for i++; i != len(newData); j++ {
				if newData[i].UserId != newData[j].UserId {
					break
				}
			}
			newData2 = append(newData2, Department{
				Employees: newData[j:i],
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
		QueryGroup(data, "UserId,Age", func(users []User) Department {
			return Department{
				Employees: users,
			}
		})
	}
}
