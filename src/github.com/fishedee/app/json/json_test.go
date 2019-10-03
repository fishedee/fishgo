package json

import (
	. "github.com/fishedee/assert"
	"reflect"
	"testing"
	"time"
)

type User struct {
	UserId     int
	Name       string
	CreateTime time.Time
}

type Users struct {
	Count int
	Data  []User
}

type Class struct {
	ClassId  int
	Name     string
	Students []User
	Score    int `json:"sss"`
	Level    int `json:"-"`
}

func TestMarshal(t *testing.T) {
	testCases := []struct {
		data interface{}
		str  string
	}{
		{
			User{1, "fish", time.Unix(0, 0)},
			`{"userId":1,"name":"fish","createTime":"1970-01-01 08:00:00"}`,
		},
		{
			Users{
				Count: 2,
				Data: []User{
					User{3, "fish", time.Unix(1, 0)},
					User{4, "cat", time.Unix(2, 0)},
				},
			},
			`{"count":2,"data":[{"userId":3,"name":"fish","createTime":"1970-01-01 08:00:01"},{"userId":4,"name":"cat","createTime":"1970-01-01 08:00:02"}]}`,
		},
		{
			Class{
				ClassId: 5,
				Name:    "class1",
				Students: []User{
					User{6, "dog", time.Unix(3, 0)},
					User{7, "apple", time.Unix(4, 0)},
				},
				Score: 78,
			},
			`{"classId":5,"name":"class1","students":[{"userId":6,"name":"dog","createTime":"1970-01-01 08:00:03"},{"userId":7,"name":"apple","createTime":"1970-01-01 08:00:04"}],"sss":78}`,
		},
	}

	//序列化
	for _, singleTestCase := range testCases {
		str, err := JsonMarshal(singleTestCase.data)
		AssertEqual(t, err, nil)
		AssertEqual(t, string(str), singleTestCase.str+"\n")
	}

	//反序列化
	for _, singleTestCase := range testCases {
		typ := reflect.TypeOf(singleTestCase.data)
		temp := reflect.New(typ).Interface()
		err := JsonUnmarshal([]byte(singleTestCase.str), temp)
		AssertEqual(t, err, nil)
		AssertEqual(t, reflect.ValueOf(temp).Elem().Interface(), singleTestCase.data)
	}
}
