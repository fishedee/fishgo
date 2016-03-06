package language

import (
	"fmt"
	"testing"
	"reflect"
)

func AssertEqual(t *testing.T,left interface{},right interface{}){
	isEqual := reflect.DeepEqual(left,right)
	if isEqual == false{
		t.Error(fmt.Sprintf("%#v != %#v",left,right))
	}
}

func TestArrayToMapBasic(t *testing.T){
	testCase := []struct{
		origin interface{}
		target interface{}
	}{
		{nil,nil},
		{true,true},
		{1,1},
		{1.2,1.2},
		{"12","12"},
		{[]int{1,2,3},[]interface{}{1,2,3}},
		{map[string]string{
			"1":"2",
			"3":"5",
		},
		map[string]interface{}{
			"1":"2",
			"3":"5",
		}},
	}
	for _,singleTestCase := range testCase{
		data := ArrayToMap(singleTestCase.origin,"json")
		AssertEqual(t,data,singleTestCase.target)
	}
}

func TestArrayToMapStruct(t *testing.T){
	testCase := []struct{
		origin interface{}
		target interface{}
	}{
		{struct{
			First string
			Second string 
			Third string `json:"Third"`
			Forth string `json:"Forth,omitempty"`
			Fifth string `json:"Fifth,omitempty"`
			Sixth string `json:"-"`
		}{"1","2","3","4","","6"},map[string]interface{}{
			"first":"1",
			"second":"2",
			"Third":"3",
			"Forth":"4",
		}},
	}
	for _,singleTestCase := range testCase{
		data := ArrayToMap(singleTestCase.origin,"json")
		AssertEqual(t,data,singleTestCase.target)
	}
}

func TestArrayToMapTotal(t *testing.T){
	origin := []struct{
		First string
		Second string 
		Third int `json:"Third"`
		Forth string `json:"Forth,omitempty"`
		Fifth string `json:"Fifth,omitempty"`
		Sixth string `json:"-"`
	}{
		{"1","2",3,"4","","6"},
		{"11","22",33,"44","55","66"},
	}
	target := []interface{}{
		map[string]interface{}{
			"first":"1",
			"second":"2",
			"Third":3,
			"Forth":"4",
		},
		map[string]interface{}{
			"first":"11",
			"second":"22",
			"Third":33,
			"Forth":"44",
			"Fifth":"55",
		},
	}
	data := ArrayToMap(origin,"json")
	AssertEqual(t,data,target)
}

func TestMapToArrayBasic(t* testing.T){
	testCase := []struct{
		origin interface{}
		target interface{}
	}{
		//basic type
		{true,true},
		{false,false},
		{"true",true},
		{"false",false},
		{-1,-1},
		{"-1",-1},
		{uint(1),uint(1)},
		{"1",uint(1)},
		{1.2,1.2},
		{"1.2",1.2},
		{true,"true"},
		{-1,"-1"},
		{uint(1),"1"},
		{1.2,"1.2"},
		{"abc","abc"},
		//array type
		{[]int{1,2,3},[]int{1,2,3}},
		{[]int{1,2,3},[]string{"1","2","3"}},
		//map type
		{map[int]int{
			1:1,
			2:2,
			3:3,
			4:4,
		},
		map[int]int{
			1:1,
			2:2,
			3:3,
			4:4,
		}},
		{map[int]int{
			1:1,
			2:2,
			3:3,
			4:4,
		},
		map[int]string{
			1:"1",
			2:"2",
			3:"3",
			4:"4",
		}},
		{map[int]int{
			1:1,
			2:2,
			3:3,
			4:4,
		},
		map[string]int{
			"1":1,
			"2":2,
			"3":3,
			"4":4,
		}},
		{map[int]int{
			1:1,
			2:2,
			3:3,
			4:4,
		},
		map[string]string{
			"1":"1",
			"2":"2",
			"3":"3",
			"4":"4",
		}},
	}
	for _,singleTestCase := range testCase{
		origin := singleTestCase.origin
		target := singleTestCase.target
		targetType := reflect.TypeOf(target)
		result := reflect.New(targetType)
		err := MapToArray(origin,result.Interface(),"json")
		AssertEqual(t,err,nil)
		AssertEqual(t,result.Elem().Interface(),target)
	}
}

func TestMapToArrayStruct(t *testing.T){
	testCase := []struct{
		origin interface{}
		target interface{}
	}{
		{map[string]interface{}{
			"first":"1",
			"second":"2",
			"Third":"3",
			"Forth":"4",
			"fifth":"5",
		},struct{
			First string
			Second int 
			Third string `json:"Third"`
			Forth string `json:"Forth,omitempty"`
			Fifth string `json:"-"`
		}{"1",2,"3","4",""}},
	}
	for _,singleTestCase := range testCase{
		origin := singleTestCase.origin
		target := singleTestCase.target
		targetType := reflect.TypeOf(target)
		result := reflect.New(targetType)
		err := MapToArray(origin,result.Interface(),"json")
		AssertEqual(t,err,nil)
		AssertEqual(t,result.Elem().Interface(),target)
	}
}

type totalTempStruct struct{
	A string
	B int
}

func TestMapToArrayTotal(t *testing.T){
	testCase := []struct{
		origin interface{}
		target interface{}
	}{
		{map[string]interface{}{
			"first":"1",
			"second":"2",
			"Third":"3",
			"Forth":"4",
			"fifth":"5",
			"sixth":[]map[string]interface{}{
				{"a":"1","b":"2"},
				{"a":"3","b":"4"},
			},
		},struct{
			First string
			Second int 
			Third string `json:"Third"`
			Forth string `json:"Forth,omitempty"`
			Fifth string `json:"-"`
			Sixth []totalTempStruct
		}{
			"1",
			2,
			"3",
			"4",
			"",
			[]totalTempStruct{
				{"1",2},
				{"3",4},
			},
		}},
	}
	for _,singleTestCase := range testCase{
		origin := singleTestCase.origin
		target := singleTestCase.target
		targetType := reflect.TypeOf(target)
		result := reflect.New(targetType)
		err := MapToArray(origin,result.Interface(),"json")
		AssertEqual(t,err,nil)
		AssertEqual(t,result.Elem().Interface(),target)
	}
}