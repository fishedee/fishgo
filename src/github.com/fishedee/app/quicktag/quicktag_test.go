package quicktag

import (
	"bytes"
	"encoding/json"
	"fmt"
	. "github.com/fishedee/assert"
	"reflect"
	"testing"
	"time"
)

var (
	jsonQuickTag *QuickTag
)

func jsonMarshal(data interface{}) ([]byte, error) {
	quickTagInstance := jsonQuickTag.GetTagInstance(data)

	buffer := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "")
	err := encoder.Encode(quickTagInstance)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func jsonUnmarshal(in []byte, data interface{}) error {
	quickTagInstance := jsonQuickTag.GetTagInstance(data)

	err := json.Unmarshal(in, quickTagInstance)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	jsonQuickTag = NewQuickTag("json")
}

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

type Article struct {
	ArticleId int
	Data      json.RawMessage
}

func TestNil(t *testing.T) {
	var data interface{}

	data = nil

	AssertEqual(t, jsonQuickTag.GetTagInstance(data), nil)
}

func TestMarshalAndUnmarshal(t *testing.T) {
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
		{
			Article{
				ArticleId: 1,
				Data:      []byte(`{"userId":3,"name":"fish"}`),
			},
			`{"articleId":1,"data":{"userId":3,"name":"fish"}}`,
		},
	}

	//序列化
	for _, singleTestCase := range testCases {
		str, err := jsonMarshal(singleTestCase.data)
		AssertEqual(t, err, nil)
		AssertEqual(t, string(str), singleTestCase.str+"\n")
	}

	//反序列化
	for _, singleTestCase := range testCases {
		typ := reflect.TypeOf(singleTestCase.data)
		temp := reflect.New(typ).Interface()
		err := jsonUnmarshal([]byte(singleTestCase.str), temp)
		fmt.Println(err)
		AssertEqual(t, err, nil)
		AssertEqual(t, reflect.ValueOf(temp).Elem().Interface(), singleTestCase.data)
	}
}
