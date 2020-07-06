package language

import (
	. "github.com/fishedee/assert"
	"testing"
	"time"
)

func testTimeSingle(t *testing.T, index int, zeroTime time.Time, shouldEqual1 bool, shouldEqual2 bool) {
	layout := "2006-01-02 15:04:05.999999999-07:00"

	//sqlite将时间Format为字符串以后保存到数据库，
	//读取数据库时，用ParseInLocation读出再用In转换时区
	zeroTimeStr := zeroTime.Format(layout)
	zeroTime2, err := time.ParseInLocation(layout, zeroTimeStr, time.UTC)
	if err != nil {
		panic(err)
	}
	zeroTime2 = zeroTime2.In(time.Local)

	AssertEqual(t, zeroTime.Equal(zeroTime2), shouldEqual1, index)

	zeroTimeStr2 := zeroTime2.Format(layout)
	zeroTime3, err := time.ParseInLocation(layout, zeroTimeStr2, time.UTC)
	if err != nil {
		panic(err)
	}
	zeroTime3 = zeroTime3.In(time.Local)

	AssertEqual(t, zeroTime.Equal(zeroTime3), shouldEqual2, index)
}
func TestAll(t *testing.T) {
	testCase := []struct {
		Data              time.Time
		ShouldEqualFirst  bool
		ShouldEqualSecond bool
	}{
		//1900及之前的都无法通过测试
		{
			time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			true,
			false,
		},
		{
			time.Date(1, 1, 1, 0, 0, 0, 0, time.Local),
			false,
			false,
		},
		{
			time.Date(1000, 1, 1, 0, 0, 0, 0, time.Local),
			false,
			false,
		},
		{
			time.Date(1900, 12, 31, 0, 0, 0, 0, time.UTC),
			true,
			false,
		},
		{
			time.Date(1900, 12, 31, 0, 0, 0, 0, time.Local),
			false,
			false,
		},
		//1901年以后才能通过测试
		{
			time.Date(1901, 1, 1, 0, 0, 0, 0, time.UTC),
			true,
			true,
		},
		{
			time.Date(1901, 1, 1, 0, 0, 0, 0, time.Local),
			true,
			true,
		},
		{
			ZERO_TIME,
			true,
			true,
		},
	}

	for index, singleCase := range testCase {
		testTimeSingle(t, index, singleCase.Data, singleCase.ShouldEqualFirst, singleCase.ShouldEqualSecond)
	}
}
