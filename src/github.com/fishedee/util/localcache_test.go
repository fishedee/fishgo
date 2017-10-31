package util

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestLocalCache(t *testing.T) {
	cache := NewLocalCache()
	task := NewTask()
	task.SetThreadCount(4)
	task.SetBufferCount(1024)
	task.SetHandler(func(data int) {
		if cache.Get(strconv.Itoa(data)) != nil {
			return
		}
		cache.Set(strconv.Itoa(data), func(id string) interface{} {
			fmt.Println(data)
			if data >= 100 {
				task.Stop()
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
