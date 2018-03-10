package queue

import (
	. "github.com/fishedee/app/log"
	. "github.com/fishedee/assert"
	"github.com/garyburd/redigo/redis"
	"github.com/streadway/amqp"
	"os/exec"
	"sort"
	"strconv"
	"testing"
	"time"
)

func newQueueForTest(t *testing.T, config QueueConfig) Queue {
	if config.Driver == "rabbitmq" {
		cmd := exec.Command("/usr/local/sbin/rabbitmqctl", "delete_vhost", "test")
		cmd.Run()
		cmd2 := exec.Command("/usr/local/sbin/rabbitmqctl", "add_vhost", "test")
		cmd2.Run()
		cmd3 := exec.Command("/usr/local/sbin/rabbitmqctl", "set_permissions", "-p", "test", "guest", ".*", ".*", ".*")
		cmd3.Run()
	}
	log, err := NewLog(LogConfig{
		Driver: "console",
	})
	if err != nil {
		panic(err)
	}
	manager, err := NewQueue(log, config)
	if err != nil {
		panic(err)
	}
	if config.Driver == "redis" {
		redisPool := manager.(*queueImplement).store.(*redisQueueStore).redisPool
		c := redisPool.Get()
		defer c.Close()
		c.Do("flushall")
	}
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
		newQueueForTest(t, QueueConfig{
			SavePath: "amqp://guest:guest@localhost:5672/test",
			Driver:   "rabbitmq",
		}),
	}

	for managerIndex, manager := range testCaseDriver {
		go manager.Run()
		for poolSize := -1; poolSize <= 2; poolSize++ {
			for singleTestCaseIndex, singleTestCase := range testCase {
				queueNameId := strconv.Itoa(singleTestCaseIndex) + "_" + strconv.Itoa(poolSize)

				//单一消费组
				topicId := "TestQueueConsume1_" + queueNameId
				manager.MustConsume(topicId, topicId+"queue", poolSize, singleTestCase.target)
				if managerIndex == 1 {
					manager.(*queueImplement).store.(*redisQueueStore).updateRouter()
				}
				manager.MustProduce(topicId, singleTestCase.origin...)
				testCaseResult := <-testCaseResultChannel
				AssertEqual(t, testCaseResult, singleTestCase.origin, singleTestCaseIndex)

				//双消费组
				topicId = "TestQueueConsume2_" + queueNameId
				manager.MustConsume(topicId, topicId+"queue", poolSize, singleTestCase.target)
				manager.MustConsume(topicId, topicId+"queue2", poolSize, singleTestCase.target)
				if managerIndex == 1 {
					manager.(*queueImplement).store.(*redisQueueStore).updateRouter()
				}
				manager.MustProduce(topicId, singleTestCase.origin...)
				testCaseResult = <-testCaseResultChannel
				AssertEqual(t, testCaseResult, singleTestCase.origin, singleTestCaseIndex)
				testCaseResult = <-testCaseResultChannel
				AssertEqual(t, testCaseResult, singleTestCase.origin, singleTestCaseIndex)

			}
		}
		manager.Close()
	}
}

func TestQueuePoolSize(t *testing.T) {
	testCaseDriver := []Queue{
		newQueueForTest(t, QueueConfig{
			SavePrefix: "queue:",
			Driver:     "memory",
		}),
	}

	testCase := []struct {
		poolSize    int
		minDuration time.Duration
		maxDuration time.Duration
	}{
		{1, time.Millisecond * 1000, time.Millisecond * 1100},
		{2, time.Millisecond * 500, time.Millisecond * 600},
	}

	for queueIndex, queue := range testCaseDriver {
		for index, test := range testCase {
			testCaseIndex := strconv.Itoa(queueIndex) + "_" + strconv.Itoa(index)
			result := make(chan int, 10)
			topicId := "queue4_" + strconv.Itoa(index)
			queue.MustConsume(topicId, "queue", test.poolSize, func(data int) {
				result <- data
				time.Sleep(time.Millisecond * 100)
			})
			for i := 0; i != 10; i++ {
				queue.MustProduce(topicId, i)
			}
			go queue.Close()
			begin := time.Now()
			queue.Run()
			end := time.Now()
			AssertEqual(t, end.Sub(begin) >= test.minDuration, true, testCaseIndex+","+end.Sub(begin).String())
			AssertEqual(t, end.Sub(begin) <= test.maxDuration, true, testCaseIndex)

			close(result)
			temp := []int{}
			for single := range result {
				temp = append(temp, single)
			}
			sort.Ints(temp)
			AssertEqual(t, temp, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, testCaseIndex)
		}

	}
}

func TestQueueClose(t *testing.T) {
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
		{newQueueForTest(t, QueueConfig{
			SavePath: "amqp://guest:guest@localhost:5672/test",
			Driver:   "rabbitmq",
		}), 789},
	}
	for singleTestCaseIndex, singleTestCase := range testCase {
		var result int
		inputEvent := make(chan bool)
		singleTestCase.Queue.MustConsume("topic", "queue", 1, func(data int) {
			inputEvent <- true
			time.Sleep(time.Second)
			result = singleTestCase.Data
		})
		if singleTestCaseIndex == 1 {
			singleTestCase.Queue.(*queueImplement).store.(*redisQueueStore).updateRouter()
		}
		singleTestCase.Queue.MustProduce("topic", singleTestCase.Data)
		<-inputEvent
		go singleTestCase.Queue.Close()
		singleTestCase.Queue.Run()
		AssertEqual(t, result, singleTestCase.Data, singleTestCaseIndex)
	}
}

func TestQueueRetry(t *testing.T) {
	testCase := []Queue{
		newQueueForTest(t, QueueConfig{
			SavePath:      "127.0.0.1:6379,100,13420693396",
			SavePrefix:    "queue3:",
			Driver:        "redis",
			RetryInterval: 2,
		}),
		newQueueForTest(t, QueueConfig{
			SavePath:      "amqp://guest:guest@localhost:5672/test",
			Driver:        "rabbitmq",
			RetryInterval: 2,
		}),
	}

	for testCaseIndex, queue := range testCase {
		result := make(chan string, 64)
		queue.MustConsume("topic1", "queue", 1, func(data string) {
			result <- data
		})
		if testCaseIndex == 0 {
			queue.(*queueImplement).store.(*redisQueueStore).updateRouter()
		}
		queue.MustProduce("topic1", "mm1")
		queue.MustProduce("topic1", "mm2")
		time.Sleep(time.Second * 1)
		if testCaseIndex == 0 {
			listener := queue.(*queueImplement).store.(*redisQueueStore).listener
			listener.Range(func(key, value interface{}) bool {
				key.(redis.Conn).Close()
				return true
			})
		} else {
			listener := queue.(*queueImplement).store.(*rabbitmqQueueStore).listener
			listener.Range(func(key, value interface{}) bool {
				key.(*amqp.Connection).Close()
				return true
			})
		}
		queue.MustProduce("topic1", "mm3")
		queue.MustProduce("topic1", "mm4")
		time.Sleep(time.Second * 2)
		queue.MustProduce("topic1", "mm5")
		queue.MustProduce("topic1", "mm6")
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
}
