package container

import (
	. "github.com/fishedee/assert"
	"testing"
)

func TestStackBasic(t *testing.T) {
	stack := NewStack()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	result := []int{3, 2, 1}
	length := stack.Len()
	for i := 0; i != length; i++ {
		AssertEqual(t, stack.Pop(), result[i])
	}
}

func TestStackFull(t *testing.T) {
	stack := NewStack()
	stack.Push(1)
	stack.Push(2)

	AssertEqual(t, stack.Top(), 2)
	AssertEqual(t, stack.Pop(), 2)
	AssertEqual(t, stack.Len(), 1)

	stack.Push(3)

	AssertEqual(t, stack.Top(), 3)
	AssertEqual(t, stack.Pop(), 3)
	AssertEqual(t, stack.Pop(), 1)
	AssertEqual(t, stack.Len(), 0)

	stack.Push(4)
	stack.Push(5)
	AssertEqual(t, stack.Len(), 2)
	AssertEqual(t, stack.Top(), 5)
	AssertEqual(t, stack.Pop(), 5)
	AssertEqual(t, stack.Len(), 1)

	AssertEqual(t, stack.Top(), 4)
	AssertEqual(t, stack.Pop(), 4)
	AssertEqual(t, stack.Len(), 0)

	//nothing
	AssertEqual(t, stack.Top(), nil)
}
