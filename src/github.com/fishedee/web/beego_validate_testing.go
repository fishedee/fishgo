package web

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/context"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

type BeegoValidateTestInterface interface {
	beegoValidateControllerInterface
	SetTesting(*testing.T)
	SetTestingMethod(method string)
	RequestReset()
}

type BeegoValidateTest struct {
	BeegoValidateController
	testingMethod string
	t             *testing.T
}

func (this *BeegoValidateTest) getBackTrace() string {
	stack := ""
	for i := 1; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		stack = stack + fmt.Sprintln(fmt.Sprintf("%s:%d", file, line))
	}
	return stack
}

func (this *BeegoValidateTest) Concurrent(number int, concurrency int, handler func()) {
	if number <= 0 {
		panic("benchmark numer is invalid")
	}
	if concurrency <= 0 {
		panic("benchmark concurrency is invalid")
	}
	singleConcurrency := number / concurrency
	if singleConcurrency <= 0 ||
		number%concurrency != 0 {
		panic("benchmark numer/concurrency is invalid")
	}

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < singleConcurrency; i++ {
				handler()
			}
		}()
	}
	wg.Wait()
}

func (this *BeegoValidateTest) Benchmark(number int, concurrency int, handler func(), testCase ...interface{}) {
	beginTime := time.Now().UnixNano()
	this.Concurrent(number, concurrency, handler)
	endTime := time.Now().UnixNano()

	totalTime := endTime - beginTime
	singleTime := totalTime / int64(number)
	singleReq := float64(number) / (float64(totalTime) / 1e9)
	if len(testCase) == 0 {
		testCase = []interface{}{""}
	}
	fmt.Printf(
		"%v:%v %v number %v concurrency %v / req, %.2freq / s\n",
		this.testingMethod,
		testCase[0],
		number,
		concurrency,
		time.Duration(singleTime).String(),
		singleReq,
	)
}

func (this *BeegoValidateTest) AssertEqual(left interface{}, right interface{}, testCase ...interface{}) {
	isEqual := reflect.DeepEqual(left, right)
	if isEqual {
		return
	}
	backtrace := this.getBackTrace()
	if len(testCase) == 0 {
		this.t.Errorf("%v : assertEqual Fail! %v != %v\n%s", this.testingMethod, left, right, backtrace)
	} else {
		this.t.Errorf("%v : assertEqual Fail! %v != %v\ntestCase: %+v\n%s", this.testingMethod, left, right, testCase, backtrace)
	}

}

func (this *BeegoValidateTest) SetTesting(t *testing.T) {
	this.t = t
}
func (this *BeegoValidateTest) SetTestingMethod(method string) {
	this.testingMethod = method
}

func (this *BeegoValidateTest) RandomInt() int {
	return rand.Int()
}

func (this *BeegoValidateTest) RandomString(length int) string {
	result := []rune{}
	for i := 0; i != length; i++ {
		var single rune
		randInt := rand.Int() % (10 + 26 + 26)
		if randInt < 10 {
			single = rune('0' + randInt)
		} else if randInt < 10+26 {
			single = rune('A' + (randInt - 10))
		} else {
			single = rune('a' + (randInt - 10 - 26))
		}
		result = append(result, single)
	}
	return string(result)
}

func (this *BeegoValidateTest) RequestReset() {
	ctx := context.NewContext()
	request, err := http.NewRequest("get", "/", bytes.NewReader([]byte("")))
	if err != nil {
		panic(err)
	}
	ctx.Reset(nil, request)
	this.SetAppContextInner(ctx)
	return ctx
}

func (this *BeegoValidateTest) RequestSetCookie() {

}

func (this *BeegoValidateTest) RequestGetCookie() {

}

func InitBeegoVaildateTest(t *testing.T, test BeegoValidateTestInterface) {
	//初始化test
	test.SetTesting(t)
	test.SetAppControllerInner(test)
	this.RequestReset()
	test.Prepare()

	isBenchTest := false
	for _, singleArgv := range os.Args {
		if strings.Index(singleArgv, "bench") != -1 {
			isBenchTest = true
		}
	}
	//遍历test，执行测试
	testType := reflect.TypeOf(test)
	testValue := reflect.ValueOf(test)
	testMethodNum := testType.NumMethod()
	for i := 0; i != testMethodNum; i++ {
		singleValueMethodType := testType.Method(i)
		if isBenchTest == false {
			if strings.HasPrefix(singleValueMethodType.Name, "Test") == false {
				continue
			}
		} else {
			if strings.HasPrefix(singleValueMethodType.Name, "Benchmark") == false ||
				singleValueMethodType.Name == "Benchmark" {
				continue
			}
		}
		test.SetTestingMethod(singleValueMethodType.Name)
		//执行测试
		singleValueMethodType.Func.Call([]reflect.Value{testValue})
	}
}

func init() {
	rand.Seed(time.Now().Unix())
}
