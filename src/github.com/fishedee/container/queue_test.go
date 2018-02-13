package container

import (
	. "github.com/fishedee/assert"
	"testing"
)

func TestQueueBasic(t *testing.T) {
	queue := NewQueue()
	queue.Push(1)
	queue.Push(2)
	queue.Push(3)

	result := []int{1, 2, 3}
	length := queue.Len()
	for i := 0; i != length; i++ {
		AssertEqual(t, queue.Pop(), result[i])
	}
}

func TestQueueFull(t *testing.T) {
	queue := NewQueue()
	queue.Push(1)
	queue.Push(2)

	AssertEqual(t, queue.Top(), 1)
	AssertEqual(t, queue.Pop(), 1)
	AssertEqual(t, queue.Len(), 1)

	queue.Push(3)

	AssertEqual(t, queue.Top(), 2)
	AssertEqual(t, queue.Pop(), 2)
	AssertEqual(t, queue.Pop(), 3)
	AssertEqual(t, queue.Len(), 0)

	queue.Push(4)
	queue.Push(5)
	AssertEqual(t, queue.Len(), 2)
	AssertEqual(t, queue.Top(), 4)
	AssertEqual(t, queue.Pop(), 4)
	AssertEqual(t, queue.Len(), 1)

	AssertEqual(t, queue.Top(), 5)
	AssertEqual(t, queue.Pop(), 5)
	AssertEqual(t, queue.Len(), 0)

	//nothing
	AssertEqual(t, queue.Top(), nil)
}
