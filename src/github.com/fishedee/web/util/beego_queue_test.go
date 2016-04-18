package util

import (
	. "github.com/fishedee/web/util/beego_queue"
	"reflect"
	"strconv"
	"testing"
)

func assertQueueEqual(t *testing.T, left interface{}, right interface{}, index int) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("case :%v ,%+v != %+v", index, left, right)
	}
}

func newQueueManagerForTest(t *testing.T, config QueueManagerConfig) *QueueManager {
	manager, err := newQueueManager(config)
	assertQueueEqual(t, err, nil, 0)
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

	testCaseDriver := []*QueueManager{
		newQueueManagerForTest(t, QueueManagerConfig{
			Driver: "memory",
		}),
		newQueueManagerForTest(t, QueueManagerConfig{
			BeegoQueueStoreConfig: BeegoQueueStoreConfig{
				SavePath: "127.0.0.1:6379,100,13420693396",
			},
			Driver: "redis",
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
	manager := newQueueManagerForTest(t, QueueManagerConfig{
		Driver: "memory",
	})
	for i := 0; i != 100; i++ {
		var result int
		manager.ConsumeInPool("queue", func(data int) {
			result = data
		}, 1)
		manager.Produce("queue", i)
		assertQueueEqual(t, i, result, i)
	}

	//ConsumeInPool配置Async
	manager2 := newQueueManagerForTest(t, QueueManagerConfig{
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
	assertQueueEqual(t, hasFalse, true, 0)

	//config配置Sync
	manager3 := newQueueManagerForTest(t, QueueManagerConfig{
		Driver:   "memory",
		PoolSize: 1,
	})
	for i := 0; i != 100; i++ {
		var result int
		manager3.Consume("queue", func(data int) {
			result = data
		})
		manager3.Produce("queue", i)
		assertQueueEqual(t, i, result, i)
	}
}
