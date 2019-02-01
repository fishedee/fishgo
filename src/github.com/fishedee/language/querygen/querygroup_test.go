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
		nextData := make([]int, len(data), len(data))
		findMap := make(map[int]int, len(data))
		newData := make([]User, len(data), len(data))
		result := make([]Department, 0, len(data))
		for i, single := range newData {
			lastIndex, isExist := findMap[single.UserId]
			if isExist {
				nextData[lastIndex] = i
			}
			nextData[i] = -1
			findMap[single.UserId] = i
		}
		for i := 0; i != len(nextData); i++ {
			j := i
			k := 0
			if nextData[j] == 0 {
				continue
			}
			for nextData[j] != -1 {
				nextJ := nextData[j]
				newData[k] = data[j]
				nextData[j] = 0
				k++
				j = nextJ
			}
			newData[k] = data[j]
			k++
			result = append(result, Department{
				Employees: newData[:k],
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
