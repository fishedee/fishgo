package queue

import (
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/assert"
	"strconv"
	"testing"
	"time"
)

func newQueueForTest(t *testing.T, config QueueConfig) Queue {
	log, err := NewLog(LogConfig{
		Driver: "console",
	})
	AssertEqual(t, err, nil, 0)
	manager, err := NewQueue(log, config)
	AssertEqual(t, err, nil, 0)
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
		{[]interface{}{}, func() {
			testCaseResultChannel <- []interface{}{}
		}},
		{[]interface{}{true}, func(data1 bool) {
			testCaseResultChannel <- []interface{}{data1}
		}},
		{[]interface{}{1}, func(data1 int) {
			testCaseResultChannel <- []interface{}{data1}
		}},
		{[]interface{}{"a"}, func(data1 string) {
			testCaseResultChannel <- []interface{}{data1}
		}},
		{[]interface{}{testStructVal}, func(data1 testStruct) {
			testCaseResultChannel <- []interface{}{data1}
		}},
		{[]interface{}{true, 1, "a", testStructVal}, func(data1 bool, data2 int, data3 string, data4 testStruct) {
			testCaseResultChannel <- []interface{}{data1, data2, data3, data4}
		}},
		//多余参数
		{[]interface{}{1, 1}, func(data1 int, data2 int) {
			testCaseResultChannel <- []interface{}{data1, data2}
		}},
		{[]interface{}{1, "aa"}, func(data1 int) {
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
				AssertEqual(t, testCaseResult, singleTestCase.origin, singleTestCaseIndex)

				//发布订阅模式，空订阅者
				queueName = "TestQueueSubscribe1" + queueNameId
				manager.Publish(queueName, singleTestCase.origin...)

				//发布订阅模式，单订阅者
				queueName = "TestQueueSubscribe2" + queueNameId
				manager.SubscribeInPool(queueName, singleTestCase.target, poolSize)
				manager.Publish(queueName, singleTestCase.origin...)
				testCaseResult = <-testCaseResultChannel
				AssertEqual(t, testCaseResult, singleTestCase.origin, singleTestCaseIndex)

				//发布订阅模式，两订阅者
				queueName = "TestQueueSubscribe3" + queueNameId
				manager.SubscribeInPool(queueName, singleTestCase.target, poolSize)
				manager.SubscribeInPool(queueName, singleTestCase.target, poolSize)
				manager.Publish(queueName, singleTestCase.origin...)
				testCaseResult = <-testCaseResultChannel
				AssertEqual(t, testCaseResult, singleTestCase.origin, singleTestCaseIndex)
				testCaseResult = <-testCaseResultChannel
				AssertEqual(t, testCaseResult, singleTestCase.origin, singleTestCaseIndex)
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
		manager.ConsumeInPool("queue", func(data int) {
			result = data
		}, 1)
		manager.Produce("queue", i)
		AssertEqual(t, i, result, i)
	}

	//ConsumeInPool配置Async
	manager2 := newQueueForTest(t, QueueConfig{
		Driver: "memory",
	})
	var hasFalse bool
	for i := 0; i != 100; i++ {
		var result int
		manager2.Consume("queue", func(data int) {
			result = data
		})
		manager2.Produce("queue", i)
		if result == 0 {
			hasFalse = true
		}
	}
	AssertEqual(t, hasFalse, true, 0)

	//config配置Sync
	manager3 := newQueueForTest(t, QueueConfig{
		Driver:   "memory",
		PoolSize: 1,
	})
	for i := 0; i != 100; i++ {
		var result int
		manager3.Consume("queue", func(data int) {
			result = data
		})
		manager3.Produce("queue", i)
		AssertEqual(t, i, result, i)
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
	queue.Consume("topic1", func(data int) {
		result <- data
	})
	target := <-result
	AssertEqual(t, target, 1, 0)

	//发布订阅模式没有累积效应
	queue.Publish("topic2", 2)
	var result2 = make(chan int)
	queue.Subscribe("topic2", func(data int) {
		result2 <- data
	})
	select {
	case <-result2:
		AssertEqual(t, false, "invalid!", 0)
	case <-time.NewTimer(time.Second).C:
		break
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
			SavePrefix: "queue2:",
			Driver:     "redis",
		}), 456},
	}
	for singleTestCaseIndex, singleTestCase := range testCase {
		var result int
		inputEvent := make(chan bool)
		singleTestCase.Queue.Consume("queue", func(data int) {
			inputEvent <- true
			time.Sleep(time.Second)
			result = singleTestCase.Data
		})
		singleTestCase.Queue.Produce("queue", singleTestCase.Data)
		<-inputEvent
		singleTestCase.Queue.Close()
		AssertEqual(t, result, singleTestCase.Data, singleTestCaseIndex)
	}
}

func TestQueueRedisRetry(t *testing.T) {
	queue := newQueueForTest(t, QueueConfig{
		SavePath:      "127.0.0.1:6379,100,13420693396",
		SavePrefix:    "queue3:",
		Driver:        "redis",
		RetryInterval: 2,
	})
	result := make(chan string, 64)
	queue.Consume("topic1", func(data string) {
		result <- data
	})
	queue.Produce("topic1", "mm1")
	queue.Produce("topic1", "mm2")
	time.Sleep(time.Second * 1)
	queue.(*queueImplement).store.(*BasicQueueStore).QueueStoreBasicInterface.(*RedisQueueStore).closeListener()
	queue.Produce("topic1", "mm3")
	queue.Produce("topic1", "mm4")
	time.Sleep(time.Second * 2)
	queue.Produce("topic1", "mm5")
	queue.Produce("topic1", "mm6")
	time.Sleep(time.Second * 2)
	testCase := []string{"mm1", "mm2", "mm3", "mm4", "mm5", "mm6"}
	for i := 0; i != 6; i++ {
		select {
		case single := <-result:
			AssertEqual(t, single, testCase[i])
		default:
			AssertEqual(t, true, false)
		}
	}
}
