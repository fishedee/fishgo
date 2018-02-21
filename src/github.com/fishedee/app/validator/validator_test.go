package validator

import (
	. "github.com/fishedee/assert"
	"io/ioutil"
	"net/http"
	"strings"
	//"reflect"
	"mime/multipart"
	"testing"
)

type testStruct struct {
	A int
	B bool
	C string
}

func TestValidatorParam(t *testing.T) {
	testCase := []struct {
		param map[string]string
		name  string
		value interface{}
	}{
		{nil, "a", ""},
		{nil, "b", ""},
		{map[string]string{
			"a": "a_value",
			"b": "b_value",
		}, "a", "a_value"},
		{map[string]string{
			"a": "a_value",
			"b": "b_value",
		}, "c", ""},
		{map[string]string{
			"a": "123",
			"b": "true",
			"c": "c_value",
		}, "", testStruct{123, true, "c_value"}},
	}

	for _, singleTestCase := range testCase {
		r, _ := http.NewRequest("GET", "http://www.baidu.com/", nil)
		validatorFactory, _ := NewValidatorFactory(ValidatorConfig{})
		validator := validatorFactory.Create(r, singleTestCase.param)

		if singleTestCase.name != "" {
			result, err := validator.Param(singleTestCase.name)
			AssertEqual(t, err, nil)
			AssertEqual(t, result, singleTestCase.value)

			result2 := validator.MustParam(singleTestCase.name)
			AssertEqual(t, result2, singleTestCase.value)
		} else {
			result := testStruct{}
			err := validator.BindParam(&result)
			AssertEqual(t, err, nil)
			AssertEqual(t, result, singleTestCase.value)

			result2 := testStruct{}
			validator.MustBindParam(&result2)
			AssertEqual(t, result, singleTestCase.value)
		}

	}
}

func TestValidatorQuery(t *testing.T) {
	testCase := []struct {
		query string
		name  string
		value interface{}
	}{
		{"", "a", ""},
		{"", "b", ""},
		{"a=a_value&b=b_value", "a", "a_value"},
		{"a=a_value&b=b_value", "c", ""},
		{"a=123&b=true&c=c_value", "", testStruct{123, true, "c_value"}},
	}

	for _, singleTestCase := range testCase {
		r, _ := http.NewRequest("GET", "http://www.baidu.com/?"+singleTestCase.query, nil)
		validatorFactory, _ := NewValidatorFactory(ValidatorConfig{})
		validator := validatorFactory.Create(r, nil)

		if singleTestCase.name != "" {
			result, err := validator.Query(singleTestCase.name)
			AssertEqual(t, err, nil)
			AssertEqual(t, result, singleTestCase.value)

			result2 := validator.MustQuery(singleTestCase.name)
			AssertEqual(t, result2, singleTestCase.value)
		} else {
			result := testStruct{}
			err := validator.BindQuery(&result)
			AssertEqual(t, err, nil)
			AssertEqual(t, result, singleTestCase.value)

			result2 := testStruct{}
			validator.MustBindQuery(&result2)
			AssertEqual(t, result, singleTestCase.value)
		}
	}
}

func TestValidatorFormPost(t *testing.T) {
	testCase := []struct {
		form  string
		name  string
		value interface{}
	}{
		{"", "a", ""},
		{"", "b", ""},
		{"a=a_value&b=b_value", "a", "a_value"},
		{"a=a_value&b=b_value", "c", ""},
		{"a=123&b=true&c=c_value", "", testStruct{123, true, "c_value"}},
	}

	for k := 0; k != 2; k++ {
		for _, singleTestCase := range testCase {
			r, _ := http.NewRequest("POST", "http://www.baidu.com/", ioutil.NopCloser(strings.NewReader(singleTestCase.form)))
			if k == 0 {
				r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
			validatorFactory, _ := NewValidatorFactory(ValidatorConfig{})
			validator := validatorFactory.Create(r, nil)

			if singleTestCase.name != "" {
				result, err := validator.Form(singleTestCase.name)
				AssertEqual(t, err, nil)
				AssertEqual(t, result, singleTestCase.value)

				result2 := validator.MustForm(singleTestCase.name)
				AssertEqual(t, result2, singleTestCase.value)
			} else {
				result := testStruct{}
				err := validator.BindForm(&result)
				AssertEqual(t, err, nil)
				AssertEqual(t, result, singleTestCase.value)

				result2 := testStruct{}
				validator.MustBindForm(&result2)
				AssertEqual(t, result, singleTestCase.value)
			}
		}
	}
}

func TestValidatorFormFileValue(t *testing.T) {
	testCase := []struct {
		form  map[string]string
		name  string
		value interface{}
	}{
		{map[string]string{}, "a", ""},
		{map[string]string{}, "b", ""},
		{map[string]string{
			"a": "a_value",
			"b": "b_value",
		}, "a", "a_value"},
		{map[string]string{
			"a": "a_value",
			"b": "b_value",
		}, "c", ""},
		{map[string]string{
			"a": "123",
			"b": "true",
			"c": "c_value",
		}, "", testStruct{123, true, "c_value"}},
	}

	for _, singleTestCase := range testCase {
		builder := strings.Builder{}
		writer := multipart.NewWriter(&builder)
		for key, value := range singleTestCase.form {
			writer.WriteField(key, value)
		}
		writer.Close()

		r, _ := http.NewRequest("POST", "http://www.baidu.com/", ioutil.NopCloser(strings.NewReader(builder.String())))
		r.Header.Set("Content-Type", writer.FormDataContentType())
		validatorFactory, _ := NewValidatorFactory(ValidatorConfig{})
		validator := validatorFactory.Create(r, nil)

		if singleTestCase.name != "" {
			result, err := validator.Form(singleTestCase.name)
			AssertEqual(t, err, nil)
			AssertEqual(t, result, singleTestCase.value)

			result2 := validator.MustForm(singleTestCase.name)
			AssertEqual(t, result2, singleTestCase.value)
		} else {
			result := testStruct{}
			err := validator.BindForm(&result)
			AssertEqual(t, err, nil)
			AssertEqual(t, result, singleTestCase.value)

			result2 := testStruct{}
			validator.MustBindForm(&result2)
			AssertEqual(t, result, singleTestCase.value)
		}

	}
}

func TestValidatorFormFile(t *testing.T) {
	testCase := []struct {
		form     map[string]string
		name     string
		filename string
		value    interface{}
	}{
		{map[string]string{
			"file1": "a.txt",
			"file2": "b.txt",
		}, "file1", "a.txt", "Hello a"},
		{map[string]string{
			"file1": "a.txt",
			"file2": "b.txt",
		}, "file2", "b.txt", "hello b"},
		{map[string]string{
			"file1": "a.txt",
			"file2": "b.txt",
		}, "file3", "", nil},
	}

	readFileData := func(header *multipart.FileHeader) string {
		file, _ := header.Open()
		data, _ := ioutil.ReadAll(file)
		return string(data)
	}
	for _, singleTestCase := range testCase {
		builder := strings.Builder{}
		writer := multipart.NewWriter(&builder)
		for key, value := range singleTestCase.form {
			fileWriter, _ := writer.CreateFormFile(key, value)
			fileData, _ := ioutil.ReadFile("testdata/" + value)
			fileWriter.Write(fileData)
		}
		writer.Close()

		r, _ := http.NewRequest("POST", "http://www.baidu.com/", ioutil.NopCloser(strings.NewReader(builder.String())))
		r.Header.Set("Content-Type", writer.FormDataContentType())
		validatorFactory, _ := NewValidatorFactory(ValidatorConfig{})
		validator := validatorFactory.Create(r, nil)

		result, err := validator.File(singleTestCase.name)
		AssertEqual(t, err, nil)
		if singleTestCase.value != nil {
			AssertEqual(t, result.Filename, singleTestCase.filename)
			AssertEqual(t, int(result.Size), len(singleTestCase.value.(string)))
			AssertEqual(t, readFileData(result), singleTestCase.value)
		} else {
			AssertEqual(t, result, (*multipart.FileHeader)(nil))
		}

		result2 := validator.MustFile(singleTestCase.name)
		if singleTestCase.value != nil {
			AssertEqual(t, result2.Filename, singleTestCase.filename)
			AssertEqual(t, int(result2.Size), len(singleTestCase.value.(string)))
			AssertEqual(t, readFileData(result2), singleTestCase.value)
		} else {
			AssertEqual(t, result2, (*multipart.FileHeader)(nil))
		}
	}
}

func TestValidatorError(t *testing.T) {
	var err error

	//Content-Type设置错误
	r, _ := http.NewRequest("POST", "http://www.baidu.com/", strings.NewReader("a=3&b=4"))
	r.Header.Set("Content-Type", "text/plain")
	validatorFactory, _ := NewValidatorFactory(ValidatorConfig{})
	validator := validatorFactory.Create(r, nil)
	_, err = validator.Form("b")
	AssertEqual(t, err != nil, true)

	//MaxSize设置错误
	r2, _ := http.NewRequest("POST", "http://www.baidu.com/", strings.NewReader("a=3&b=4"))
	validatorFactory2, _ := NewValidatorFactory(ValidatorConfig{
		MaxBodySize: 1,
	})
	validator2 := validatorFactory2.Create(r2, nil)
	_, err = validator2.Query("b")
	AssertEqual(t, err != nil, true)

	//类型错误
	r3, _ := http.NewRequest("GET", "http://www.baidu.com/?a=123c", nil)
	validatorFactory3, _ := NewValidatorFactory(ValidatorConfig{})
	validator3 := validatorFactory3.Create(r3, nil)
	var data testStruct
	err = validator3.BindQuery(&data)
	AssertEqual(t, err != nil, true)

}
