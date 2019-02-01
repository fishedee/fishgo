package main

import (
	"fmt"
	. "github.com/fishedee/assert"
	. "github.com/fishedee/language"
	. "github.com/fishedee/language/querygen/testdata"
	"os"
	"testing"
)

func TestQueryColumn(t *testing.T) {
	data := []User{
		User{UserId: 1},
		User{UserId: -2},
		User{UserId: 3},
	}
	AssertEqual(t, QueryColumn(data, "UserId"), []int{1, -2, 3})
	AssertEqual(t, QueryColumn(data, "."), data)
	AssertEqual(t, QueryColumn([]int{1, -2, 3}, "."), []int{1, -2, 3})
}

func BenchmarkQueryColumnHand(b *testing.B) {
	data := make([]User, 1000, 1000)

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		newData := make([]int, len(data), len(data))
		for i, single := range data {
			newData[i] = single.UserId
		}
	}
}

func BenchmarkQueryColumnMacro(b *testing.B) {
	data := make([]User, 1000, 1000)

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		QueryColumn(data, "UserId")
	}
}

func BenchmarkQueryColumnReflect(b *testing.B) {
	data := make([]User, 1000, 1000)

	b.ResetTimer()
	for i := 0; i != b.N; i++ {
		QueryColumn(data, "Age")
	}
}

func init() {
	args := os.Args
	isWarning := true
	for _, arg := range args {
		if arg == "-test.benchmem=true" {
			isWarning = false
			break
		}
	}
	fmt.Println("QueryReflectWarning:", isWarning)
	QueryReflectWarning(isWarning)
}
