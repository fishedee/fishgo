package web

import (
	"fmt"
	. "github.com/fishedee/assert"
	"testing"
)

func TestOldestStayElem(t *testing.T) {
	container := NewOldestStayContainer()

	for i := 0; i != 20; i++ {
		container.Push(int64(i), fmt.Sprintf("fish_%d", i))
	}

	for i := 10; i >= 0; i-- {
		container.Pop(int64(i))
	}

	container.Pop(13)
	container.Pop(15)

	stays := container.OldestStay(5)
	t.Logf("%v", stays)
	for i, _ := range stays {
		stays[i].Timestamp = 0
	}
	AssertEqual(t, stays, []OldestStayElem{
		OldestStayElem{0, 11, "fish_11"},
		OldestStayElem{0, 12, "fish_12"},
		OldestStayElem{0, 14, "fish_14"},
		OldestStayElem{0, 16, "fish_16"},
		OldestStayElem{0, 17, "fish_17"},
	})
	for i := 21; i != 25; i++ {
		container.Push(int64(i), fmt.Sprintf("cat_%d", i))
	}

	container.Pop(12)
	container.Pop(14)
	container.Pop(18)

	stays2 := container.OldestStay(5)
	t.Logf("%v", stays2)
	for i, _ := range stays2 {
		stays2[i].Timestamp = 0
	}
	AssertEqual(t, stays2, []OldestStayElem{
		OldestStayElem{0, 11, "fish_11"},
		OldestStayElem{0, 16, "fish_16"},
		OldestStayElem{0, 17, "fish_17"},
		OldestStayElem{0, 19, "fish_19"},
		OldestStayElem{0, 21, "cat_21"},
	})

	stays4 := container.OldestStay(3)
	t.Logf("%v", stays4)
	for i, _ := range stays4 {
		stays4[i].Timestamp = 0
	}
	AssertEqual(t, stays4, []OldestStayElem{
		OldestStayElem{0, 11, "fish_11"},
		OldestStayElem{0, 16, "fish_16"},
		OldestStayElem{0, 17, "fish_17"},
	})

	container.Pop(11)
	container.Pop(17)
	container.Pop(19)
	container.Pop(21)
	container.Pop(22)
	container.Pop(23)

	stays3 := container.OldestStay(5)
	t.Logf("%v", stays3)
	for i, _ := range stays3 {
		stays3[i].Timestamp = 0
	}
	AssertEqual(t, stays3, []OldestStayElem{
		OldestStayElem{0, 16, "fish_16"},
		OldestStayElem{0, 24, "cat_24"},
	})
}
