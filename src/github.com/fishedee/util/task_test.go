package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestTask(t *testing.T) {
	task := NewTask()
	task.SetThreadCount(4)
	task.SetBufferCount(1024)
	task.SetHandler(func(data int) {
		fmt.Println(data)
		if data == 100 {
			task.Stop()
		}
	})
	task.Start()
	for i := 1; i <= 100; i++ {
		task.AddTask(i)
	}
	task.Wait()
}

func TestTask2(t *testing.T) {
	task := NewTask()
	task.SetThreadCount(2)
	task.SetBufferCount(1)
	task.SetHandler(func(data int) {
		fmt.Println(data)
		if data == 10 {
			task.Stop()
		}
		for i := 1; i <= 5; i++ {
			task.AddTask(i + data)
		}
	})
	task.Start()
	task.AddTask(1)
	task.Wait()
}

func TestTask3(t *testing.T) {
	cache := NewLocalCache()
	task := NewTask()
	task.SetIsAutoStop(true)
	task.SetThreadCount(4)
	task.SetBufferCount(1024)
	task.SetHandler(func(data int) {
		if cache.Get(strconv.Itoa(data)) != nil {
			return
		}
		cache.Set(strconv.Itoa(data), func(id string) interface{} {
			fmt.Println(data)
			if data >= 100 {
				return true
			}
			nextCount := int(rand.Uint32()%5 + 1)
			for i := 1; i <= nextCount; i++ {
				task.AddTask(i + data)
			}
			return true
		})

	})
	task.Start()
	task.AddTask(1)
	task.Wait()
}
