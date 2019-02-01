package main

import (
	. "github.com/fishedee/assert"
	. "github.com/fishedee/language"
	. "github.com/fishedee/language/querygen/testdata"
	"testing"
)

func TestQueryWhere(t *testing.T) {
	data := []User{
		User{Name: "Man_a"},
		User{Name: "Woman_b"},
		User{Name: "Man_c"},
	}

	AssertEqual(t, QueryWhere(data, func(a User) bool {
		if len(a.Name) >= 3 && a.Name[0:3] == "Man" {
			return true
		} else {
			return false
		}
	}), []User{
		User{Name: "Man_a"},
		User{Name: "Man_c"},
	})
	AssertEqual(t, QueryWhere([]int{3, 2, 3, 5, 9, 4}, func(c int) bool {
		return c%2 == 0
	}), []int{2, 4})
}

func BenchmarkQueryWhereHand(b *testing.B) {
	data := make([]User, 1000, 1000)

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		newData := make([]User, 0, len(data))
		for _, single := range data {
			isMan := func(a User) bool {
				if len(a.Name) >= 3 && a.Name[0:3] == "Man" {
					return true
				} else {
					return false
				}
			}(single)

			if isMan {
				newData = append(newData, single)
			}
		}
	}
}

func BenchmarkQueryWhereMacro(b *testing.B) {
	data := make([]User, 1000, 1000)

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		QueryWhere(data, func(a User) bool {
			if len(a.Name) >= 3 && a.Name[0:3] == "Man" {
				return true
			} else {
				return false
			}
		})
	}
}

func BenchmarkQueryWhereReflect(b *testing.B) {
	data := make([]Sex, 1000, 1000)

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		QueryWhere(data, func(a Sex) bool {
			if a.IsMale == true {
				return true
			} else {
				return false
			}
		})
	}
}
