package container

import (
	. "github.com/fishedee/assert"
	"math/rand"
	"strconv"
	"testing"
)

func TestHashListBasic(t *testing.T) {
	testData := []int{1, 3, 7, 11, 19, 54, 107}
	testNotFoundData := []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 22, 100}

	hashList := NewHashList(10)

	//普通的设置查找
	for _, data := range testData {
		hashList.Set(data, strconv.Itoa(data)+"_value")
	}
	for _, data := range testData {
		value := hashList.Get(data)
		AssertEqual(t, value, strconv.Itoa(data)+"_value")
	}
	for _, data := range testNotFoundData {
		value := hashList.Get(data)
		AssertEqual(t, value, nil)
	}
	AssertEqual(t, hashList.Len(), len(testData))

	//重复设置
	for index, data := range testData {
		if index >= 3 {
			continue
		}
		hashList.Set(data, strconv.Itoa(data)+"_value2")
	}
	for index, data := range testData {
		value := hashList.Get(data)
		if index >= 3 {
			AssertEqual(t, value, strconv.Itoa(data)+"_value")
		} else {
			AssertEqual(t, value, strconv.Itoa(data)+"_value2")
		}
	}
	for _, data := range testNotFoundData {
		value := hashList.Get(data)
		AssertEqual(t, value, nil)
	}
	AssertEqual(t, hashList.Len(), len(testData))

	//转换查找
	hashListArray := hashList.ToHashListArray()
	for index, data := range testData {
		value := hashListArray.Get(data)
		if index >= 3 {
			AssertEqual(t, value, strconv.Itoa(data)+"_value")
		} else {
			AssertEqual(t, value, strconv.Itoa(data)+"_value2")
		}
	}
	for _, data := range testNotFoundData {
		value := hashListArray.Get(data)
		AssertEqual(t, value, nil)
	}
	AssertEqual(t, hashListArray.Len(), len(testData))

	//删除
	for index, data := range testData {
		if index >= 5 {
			continue
		}
		hashList.Del(data)
	}
	for index, data := range testData {
		value := hashList.Get(data)
		if index >= 5 {
			AssertEqual(t, value, strconv.Itoa(data)+"_value")
		} else {
			AssertEqual(t, value, nil)
		}
	}
	for _, data := range testNotFoundData {
		value := hashList.Get(data)
		AssertEqual(t, value, nil)
	}
	AssertEqual(t, hashList.Len(), len(testData)-5)
}

func getHashData(size int) []int {
	data := []int{}
	for i := 0; i != size; i++ {
		data = append(data, rand.Int())
	}
	return data
}

func BenchmarkHashIntSpeed(b *testing.B) {
	insertData := getHashData(100)
	findData := getHashData(b.N)
	hash := NewHashList(100)
	for _, singleData := range insertData {
		hash.Set(singleData, 0)
	}
	hashArray := hash.ToHashListArray()

	b.ResetTimer()
	b.StartTimer()
	for _, singleData := range findData {
		hashArray.Get(singleData)
	}
	b.StopTimer()
}

func BenchmarkMapIntSpeed(b *testing.B) {
	insertData := getHashData(100)
	findData := getHashData(b.N)
	mapper := map[int]int{}
	for _, singleData := range insertData {
		mapper[singleData] = 0
	}

	b.ResetTimer()
	b.StartTimer()
	for _, singleData := range findData {
		_, _ = mapper[singleData]
	}
	b.StopTimer()
}
