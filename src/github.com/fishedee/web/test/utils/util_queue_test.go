package web

import (
	. "github.com/fishedee/web"
	"net/http"
	"reflect"
	"strconv"
	"testing"
	"time"
)

type queueResponseWriter struct {
	header     http.Header
	headerCode int
	data       []byte
}

func (this *queueResponseWriter) Header() http.Header {
	if this.header == nil {
		this.header = http.Header{}
	}
	return this.header
}

func (this *queueResponseWriter) Write(in []byte) (int, error) {
	this.data = append(this.data, in...)
	return len(this.data), nil
}

func (this *queueResponseWriter) WriteHeader(headerCode int) {
	this.headerCode = headerCode
}

type queueModel struct {
	Model
}

func assertQueueEqual(t *testing.T, left interface{}, right interface{}, index int) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("case :%v ,%+v != %+v", index, left, right)
	}
}

func newQueueForTest(t *testing.T, config QueueConfig) Queue {
	request, err := http.NewRequest("GET", "http://www.baidu.com", nil)
	assertQueueEqual(t, err, nil, 0)
	ctx := NewContext(request, &queueResponseWriter{}, nil)
	log, err := NewLog(LogConfig{
		Driver: "console",
	})
	assertQueueEqual(t, err, nil, 0)
	manager, err := NewQueue(config)
	assertQueueEqual(t, err, nil, 0)
	manager = manager.WithLogAndContext(log, ctx)
	return manager
}

func TestQueueBasic(t *testing.T) {
	type testStruct struct {
		TntVal   bool
		StrVal   string
		SliceVal []int
		MapVal   map[string]int
	}
	testStructVal := testStruct{
		true,
		"a",
		[]int{1, 2, 3, 4},
		map[string]int{
			"aa": 11,
			"bb": 22,
		},
	}
	testCaseResultChannel := make(chan []interface{}, 10)
	testCase := []struct {
		origin []interface{}
		target interface{}
	}{
		//基础用例
		{[]interface{}{}, func(this *queueModel) {
			testCaseResultChannel <- []interface{}{}
		}},
		{[]interface{}{true}, func(this *queueModel, data1 bool) {
			testCaseResultChannel <- []interface{}{data1}
		}},
		{[]interface{}{1}, func(this *queueModel, data1 int) {
			testCaseResultChannel <- []interface{}{data1}
		}},
		{[]interface{}{"a"}, func(this *queueModel, data1 string) {
			testCaseResultChannel <- []interface{}{data1}
		}},
		{[]interface{}{testStructVal}, func(this *queueModel, data1 testStruct) {
			testCaseResultChannel <- []interface{}{data1}
		}},
		{[]interface{}{true, 1, "a", testStructVal}, func(this *queueModel, data1 bool, data2 int, data3 string, data4 testStruct) {
			testCaseResultChannel <- []interface{}{data1, data2, data3, data4}
		}},
		//多余参数
		{[]interface{}{1, 1}, func(this *queueModel, data1 int, data2 int) {
			testCaseResultChannel <- []interface{}{data1, data2}
		}},
		{[]interface{}{1, "aa"}, func(this *queueModel, data1 int) {
			testCaseResultChannel <- []interface{}{data1, "aa"}
		}},
	}

	testCaseDriver := []Queue{
		newQueueForTest(t, QueueConfig{
			SavePrefix: "queue:",
			Driver:     "memory",
		}),
		newQueueForTest(t, QueueConfig{
			SavePath:   "127.0.0.1:6379,100,13420693396",
			SavePrefix: "queue:",
			Driver:     "redis",
		}),
	}

	for _, manager := range testCaseDriver {
		for poolSize := -1; poolSize <= 2; poolSize++ {
			for singleTestCaseIndex, singleTestCase := range testCase {
				queueNameId := strconv.Itoa(singleTestCaseIndex) + "_" + strconv.Itoa(poolSize)

				//生产者消费者模式，空消费者
				queueName := "TestQueueConsume1" + queueNameId
				manager.Produce(queueName, singleTestCase.origin...)

				//生产者消费者模式，单消费者
				queueName = "TestQueueConsume2" + queueNameId
				manager.ConsumeInPool(queueName, singleTestCase.target, poolSize)
				manager.Produce(queueName, singleTestCase.origin...)
				testCaseResult := <-testCaseResultChannel
				assertQueueEqual(t, testCaseResult, singleTestCase.origin, singleTestCaseIndex)

				//发布订阅模式，空订阅者
				queueName = "TestQueueSubscribe1" + queueNameId
				manager.Publish(queueName, singleTestCase.origin...)

				//发布订阅模式，单订阅者
				queueName = "TestQueueSubscribe2" + queueNameId
				manager.SubscribeInPool(queueName, singleTestCase.target, poolSize)
				manager.Publish(queueName, singleTestCase.origin...)
				testCaseResult = <-testCaseResultChannel
				assertQueueEqual(t, testCaseResult, singleTestCase.origin, singleTestCaseIndex)

				//发布订阅模式，两订阅者
				queueName = "TestQueueSubscribe3" + queueNameId
				manager.SubscribeInPool(queueName, singleTestCase.target, poolSize)
				manager.SubscribeInPool(queueName, singleTestCase.target, poolSize)
				manager.Publish(queueName, singleTestCase.origin...)
				testCaseResult = <-testCaseResultChannel
				assertQueueEqual(t, testCaseResult, singleTestCase.origin, singleTestCaseIndex)
				testCaseResult = <-testCaseResultChannel
				assertQueueEqual(t, testCaseResult, singleTestCase.origin, singleTestCaseIndex)
			}
		}
	}
}

func TestQueueSync(t *testing.T) {
	//ConsumeInPool配置Sync
	manager := newQueueForTest(t, QueueConfig{
		Driver: "memory",
	})
	for i := 0; i != 100; i++ {
		var result int
		manager.ConsumeInPool("queue", func(this *queueModel, data int) {
			result = data
		}, 1)
		manager.Produce("queue", i)
		assertQueueEqual(t, i, result, i)
	}

	//ConsumeInPool配置Async
	manager2 := newQueueForTest(t, QueueConfig{
		Driver: "memory",
	})
	var hasFalse bool
	for i := 0; i != 100; i++ {
		var result int
		manager2.Consume("queue", func(this *queueModel, data int) {
			result = data
		})
		manager2.Produce("queue", i)
		if result == 0 {
			hasFalse = true
		}
	}
	assertQueueEqual(t, hasFalse, true, 0)

	//config配置Sync
	manager3 := newQueueForTest(t, QueueConfig{
		Driver:   "memory",
		PoolSize: 1,
	})
	for i := 0; i != 100; i++ {
		var result int
		manager3.Consume("queue", func(this *queueModel, data int) {
			result = data
		})
		manager3.Produce("queue", i)
		assertQueueEqual(t, i, result, i)
	}
}

func TestQueueRedis(t *testing.T) {
	queue := newQueueForTest(t, QueueConfig{
		SavePath:   "127.0.0.1:6379,100,13420693396",
		SavePrefix: "queue:",
		Driver:     "redis",
	})

	//生产者消费者模式有累积效应
	queue.Produce("topic1", 1)
	var result = make(chan int)
	queue.Consume("topic1", func(this *queueModel, data int) {
		result <- data
	})
	target := <-result
	assertQueueEqual(t, target, 1, 0)

	//发布订阅模式没有累积效应
	queue.Publish("topic2", 2)
	var result2 = make(chan int)
	queue.Subscribe("topic2", func(this *queueModel, data int) {
		result2 <- data
	})
	select {
	case <-result2:
		assertQueueEqual(t, false, "invalid!", 0)
	case <-time.NewTimer(time.Second).C:
		break
	}
}

func TestQueueCtx(t *testing.T) {
	testCase := []struct {
		method string
		url    string
		header http.Header
	}{
		{"GET", "http://www.baidu.com", map[string][]string{
			"User-Agent": []string{"userAgent1"},
			"Cookie":     []string{"123"},
		}},
		{"GET", "http://www.baidu.com?a=1&b=3", map[string][]string{
			"UserAgent": []string{"userAgent1", "userAgent2"},
			"Cookie":    []string{"123"},
		}},
		{"POST", "http://www.baidu.com", map[string][]string{
			"UserAgent": []string{"userAgent1", "userAgent2"},
			"Cookie":    []string{"123"},
		}},
		{"POST", "http://www.baidu.com?a=1&b=3", map[string][]string{
			"UserAgent": []string{"userAgent1", "userAgent2"},
			"Cookie":    []string{"123"},
		}},
	}

	for singleTestCaseIndex, singleTestCase := range testCase {
		request, err := http.NewRequest(singleTestCase.method, singleTestCase.url, nil)
		request.Header = singleTestCase.header
		assertQueueEqual(t, err, nil, 0)
		ctx := NewContext(request, &queueResponseWriter{}, nil)
		log, err := NewLog(LogConfig{
			Driver: "console",
		})
		assertQueueEqual(t, err, nil, 0)
		manager, err := NewQueue(QueueConfig{
			SavePrefix: "queue:",
			Driver:     "memory",
		})
		assertQueueEqual(t, err, nil, 0)
		manager = manager.WithLogAndContext(log, ctx)

		var result *http.Request
		manager.ConsumeInPool("queue", func(this *queueModel) {
			result = this.Ctx.GetRawRequest().(*http.Request)
		}, 1)
		manager.Produce("queue")
		assertQueueEqual(t, result.Method, singleTestCase.method, singleTestCaseIndex)
		assertQueueEqual(t, result.URL.String(), singleTestCase.url, singleTestCaseIndex)
		assertQueueEqual(t, result.Header, singleTestCase.header, singleTestCaseIndex)
	}
}

func TestQueueClose(t *testing.T) {
	//ConsumeInPool配置Sync
	testCase := []struct {
		Queue Queue
		Data  int
	}{
		{newQueueForTest(t, QueueConfig{
			Driver: "memory",
		}), 123},
		{newQueueForTest(t, QueueConfig{
			SavePath:   "127.0.0.1:6379,100,13420693396",
			SavePrefix: "queue:",
			Driver:     "redis",
		}), 456},
	}
	for singleTestCaseIndex, singleTestCase := range testCase {
		var result int
		inputEvent := make(chan bool)
		singleTestCase.Queue.Consume("queue", func(this *queueModel, data int) {
			inputEvent <- true
			time.Sleep(time.Second)
			result = singleTestCase.Data
		})
		singleTestCase.Queue.Produce("queue", singleTestCase.Data)
		<-inputEvent
		singleTestCase.Queue.Close()
		assertQueueEqual(t, result, singleTestCase.Data, singleTestCaseIndex)
	}
}
