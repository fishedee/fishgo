package web

import (
	"bytes"
	"fmt"
	_ "github.com/a"
	"github.com/astaxie/beego/context"
	. "github.com/fishedee/language"
	"math/rand"
	"net/http"
	"os"
	"path"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type BeegoValidateTestInterface interface {
	beegoValidateControllerInterface
	RequestReset()
}

type BeegoValidateTestResponseWriter struct {
	header     http.Header
	headerCode int
	data       []byte
}

func (this *BeegoValidateTestResponseWriter) Header() http.Header {
	if this.header == nil {
		this.header = http.Header{}
	}
	return this.header
}

func (this *BeegoValidateTestResponseWriter) Write(in []byte) (int, error) {
	this.data = append(this.data, in...)
	return len(this.data), nil
}

func (this *BeegoValidateTestResponseWriter) WriteHeader(headerCode int) {
	this.headerCode = headerCode
}

type BeegoValidateTest struct {
	BeegoValidateController
	testingMethod string
}

func (this *BeegoValidateTest) getTraceLineNumber(traceNumber int) string {
	_, filename, line, ok := runtime.Caller(traceNumber + 1)
	if !ok {
		return "???.go:???"
	} else {
		return path.Base(filename) + ":" + strconv.Itoa(line)
	}
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
			runtime.LockOSThread()
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
	traceInfo := this.getTraceLineNumber(1)
	if len(testCase) == 0 {
		this.t.Errorf("%v: assertEqual Fail! %v != %v", traceInfo, left, right)
	} else {
		this.t.Errorf("%v:%v: assertEqual Fail! %v != %v", traceInfo, testCase[0], left, right)
	}
}

func (this *BeegoValidateTest) AssertError(left Exception, rightCode int, rightMessage string, testCase ...interface{}) {
	errorString := ""
	if left.GetCode() != rightCode {
		errorString = fmt.Sprintf("assertError Code Fail! %v != %v ", left.GetCode(), rightCode)
	}
	if left.GetMessage() != rightMessage {
		errorString = fmt.Sprintf("assertError Message Fail! %v != %v ", left.GetMessage(), rightMessage)
	}
	if errorString == "" {
		return
	}
	traceInfo := this.getTraceLineNumber(1)
	if len(testCase) == 0 {
		this.t.Errorf("%v: %v", traceInfo, errorString)
	} else {
		this.t.Errorf("%v:%v: %v", traceInfo, testCase[0], errorString)
	}

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
	ctx.Reset(&BeegoValidateTestResponseWriter{}, request)
	this.SetAppContextInner(ctx)
	this.Prepare()
}

var beegoValidateTestMap map[string][]BeegoValidateTestInterface
var beegoValidateTestMapInit bool

func init() {
	beegoValidateTestMap = map[string][]BeegoValidateTestInterface{}
	beegoValidateTestMapInit = false
}

func runBeegoValidateSingleTest(t *testing.T, test BeegoValidateTestInterface) {
	//初始化test
	test.SetAppTestInner(t)
	test.SetAppControllerInner(test)
	test.RequestReset()

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
		//执行测试
		singleValueMethodType.Func.Call([]reflect.Value{testValue})
	}
}

func RunBeegoValidateTest(t *testing.T, data interface{}) {
	//获取package
	pkgPath := reflect.TypeOf(data).Elem().PkgPath()

	//初始化runtime
	if beegoValidateTestMapInit == false {
		runtime.GOMAXPROCS(runtime.NumCPU() * 4)
		rand.Seed(time.Now().Unix())
		beegoValidateTestMapInit = true
	}

	//遍历测试
	for _, singleTest := range beegoValidateTestMap[pkgPath] {
		runBeegoValidateSingleTest(t, singleTest)
	}
}

func InitBeegoVaildateTest(test BeegoValidateTestInterface) {
	pkgPath := reflect.TypeOf(test).Elem().PkgPath()
	beegoValidateTestMap[pkgPath] = append(beegoValidateTestMap[pkgPath], test)
}
