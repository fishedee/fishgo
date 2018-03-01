package util

import (
	. "github.com/fishedee/assert"
	"testing"
	"time"
)

func TestInfiniteChannelNormal(t *testing.T) {
	channel := NewInfiniteChannel()
	for i := 0; i != 10; i++ {
		channel.Write(i)
	}
	channel.Close()
	data := []int{}
	for single := range channel.Read() {
		data = append(data, single.(int))
	}
	AssertEqual(t, data, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}

func TestInfiniteChannelReadFirst(t *testing.T) {
	channel := NewInfiniteChannel()
	data := []int{}
	go func() {
		for single := range channel.Read() {
			data = append(data, single.(int))
			time.Sleep(time.Millisecond * 100)
		}
	}()

	begin := time.Now()
	for i := 0; i != 10; i++ {
		channel.Write(i)
	}
	channel.Close()
	end := time.Now()
	duration := end.Sub(begin)

	AssertEqual(t, duration < time.Millisecond*10, true)
	time.Sleep(time.Millisecond * 1100)
	AssertEqual(t, data, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
}
