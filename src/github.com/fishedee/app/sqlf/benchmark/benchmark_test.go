package main

import (
	"fmt"
	. "github.com/fishedee/language"
	"testing"
)

type DbDriver interface {
	Init()
	GetAllMaterial() []Material
	GetProduct(productIds []int) []Product
}

func runTest(db DbDriver) int {
	materials := db.GetAllMaterial()
	productIds := QueryColumn(materials, "ProductId").([]int)
	products := db.GetProduct(productIds)
	return len(products)
}

func TestAll(t *testing.T) {
	drivers := []DbDriver{
		&PureDb{},
		&XormDb{},
		&DbrDb{},
		&SqlxDb{},
		&SqlfDb{},
	}

	for _, driver := range drivers {
		driver.Init()
		productLen := runTest(driver)
		fmt.Printf("all Products len %v\n", productLen)
	}
}

func BenchmarkPureDb(b *testing.B) {
	db := &PureDb{}
	db.Init()

	b.ResetTimer()

	for i := 0; i != b.N; i++ {
		runTest(db)
	}
}

func BenchmarkXormDb(b *testing.B) {
	db := &XormDb{}
	db.Init()

	b.ResetTimer()

	for i := 0; i != b.N; i++ {
		runTest(db)
	}
}

func BenchmarkDbrDb(b *testing.B) {
	db := &DbrDb{}
	db.Init()

	b.ResetTimer()

	for i := 0; i != b.N; i++ {
		runTest(db)
	}
}

func BenchmarkSqlxDb(b *testing.B) {
	db := &SqlxDb{}
	db.Init()

	b.ResetTimer()

	for i := 0; i != b.N; i++ {
		runTest(db)
	}
}

func BenchmarkSqlfDb(b *testing.B) {
	db := &SqlfDb{}
	db.Init()

	b.ResetTimer()

	for i := 0; i != b.N; i++ {
		runTest(db)
	}
}
