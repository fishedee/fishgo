package main

import (
	. "github.com/fishedee/assert"
	. "github.com/fishedee/language"
	. "github.com/fishedee/language/querygen/testdata"
	"math/rand"
	"testing"
)

func TestQueryColumnMap(t *testing.T) {
	data := []User{
		User{UserId: 1},
		User{UserId: -2},
		User{UserId: 3},
	}
	AssertEqual(t, QueryColumnMap(data, "UserId"), map[int]User{
		1:  User{UserId: 1},
		-2: User{UserId: -2},
		3:  User{UserId: 3},
	})
	AssertEqual(t, QueryColumnMap(data, "[]UserId"), map[int][]User{
		1:  []User{User{UserId: 1}},
		-2: []User{User{UserId: -2}},
		3:  []User{User{UserId: 3}},
	})
	AssertEqual(t, QueryColumnMap([]int{5, 6, 8, 8, 0, 6}, "."), map[int]int{
		5: 5,
		6: 6,
		8: 8,
		0: 0,
	})
}

func initQueryColumnMapData() []User {
	data := make([]User, 1000, 1000)
	for i := range data {
		data[i].UserId = rand.Int()
		data[i].Age = rand.Int()
	}
	return data
}

func BenchmarkQueryColumnMapHand(b *testing.B) {
	data := initQueryColumnMapData()

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		newData := make(map[int]User, len(data))
		for _, single := range data {
			newData[single.UserId] = single
		}
	}
}

func BenchmarkQueryColumnMapMacro(b *testing.B) {
	data := initQueryColumnMapData()

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		QueryColumnMap(data, "UserId")
	}
}

func BenchmarkQueryColumnMapReflect(b *testing.B) {
	data := initQueryColumnMapData()

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		QueryColumnMap(data, "Age")
	}
}

func BenchmarkQueryColumnMapSliceHand(b *testing.B) {
	data := initQueryColumnMapData()

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		newData := make(map[int][]User, len(data))
		for _, single := range data {
			temp := newData[single.UserId]
			temp = append(temp, single)
			newData[single.UserId] = temp
		}
	}
}

func BenchmarkQueryColumnMapSliceMacro(b *testing.B) {
	data := initQueryColumnMapData()

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		QueryColumnMap(data, "[]UserId")
	}
}

func BenchmarkQueryColumnMapSliceReflect(b *testing.B) {
	data := initQueryColumnMapData()

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		QueryColumnMap(data, "[]Age")
	}
}
