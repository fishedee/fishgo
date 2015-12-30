package util

import (
	"errors"
	"github.com/astaxie/beego"
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
	. "github.com/fishedee/web/util/beego_queue"
)

type QueueManagerConfig struct {
	Driver string `json:driver`
}

type QueueManager struct {
	BeegoQueueStoreInterface
	Log     *LogManager
	Monitor *MonitorManager
}

var newQueueManagerMemory *MemoryFunc
var newQueueManagerFromConfigMemory *MemoryFunc

func init() {
	var err error
	newQueueManagerMemory, err = NewMemoryFunc(newQueueManager, MemoryFuncCacheNormal)
	if err != nil {
		panic(err)
	}
	newQueueManagerFromConfigMemory, err = NewMemoryFunc(newQueueManagerFromConfig, MemoryFuncCacheNormal)
	if err != nil {
		panic(err)
	}
}

func newQueueManager(config QueueManagerConfig) (*QueueManager, error) {
	if config.Driver == "" {
		return nil, nil
	} else if config.Driver == "memory" {
		queue, err := NewMemoryQueue()
		if err != nil {
			return nil, err
		}
		return &QueueManager{
			BeegoQueueStoreInterface: queue,
		}, nil
	} else if config.Driver == "memory_async" {
		queue, err := NewMemoryAsyncQueue()
		if err != nil {
			return nil, err
		}
		return &QueueManager{
			BeegoQueueStoreInterface: queue,
		}, nil
	} else {
		return nil, errors.New("invalid memory config " + config.Driver)
	}
}

func NewQueueManager(config QueueManagerConfig) (*QueueManager, error) {
	result, err := newQueueManagerMemory.Call(config)
	if err != nil {
		return nil, err
	}
	return result.(*QueueManager), err
}

func newQueueManagerFromConfig(configName string) (*QueueManager, error) {
	driver := beego.AppConfig.String(configName + "driver")

	queueConfig := QueueManagerConfig{}
	queueConfig.Driver = driver
	return NewQueueManager(queueConfig)
}

func NewQueueManagerFromConfig(configName string) (*QueueManager, error) {
	result, err := newQueueManagerFromConfigMemory.Call(configName)
	if err != nil {
		return nil, err
	}
	return result.(*QueueManager), err
}

func NewQueueManagerWithLogAndMonitor(log *LogManager, monitor *MonitorManager, queue *QueueManager) *QueueManager {
	if queue == nil {
		return nil
	} else {
		return &QueueManager{
			BeegoQueueStoreInterface: queue.BeegoQueueStoreInterface,
			Log:     log,
			Monitor: monitor,
		}
	}
}

func (this *QueueManager) WrapListener(listener BeegoQueueListener) BeegoQueueListener {
	return func(argv interface{}) {
		defer CatchCrash(func(exception Exception) {
			this.Log.Critical("QueueTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
			if this.Monitor != nil {
				this.Monitor.AscCriticalCount()
			}
		})
		defer Catch(func(exception Exception) {
			this.Log.Error("QueueTask Error Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
			if this.Monitor != nil {
				this.Monitor.AscErrorCount()
			}
		})
		listener(argv)
	}
}

func (this *QueueManager) Consume(topicId string, listener BeegoQueueListener) error {
	return this.BeegoQueueStoreInterface.Consume(topicId, this.WrapListener(listener))
}

func (this *QueueManager) Subscribe(topicId string, listener BeegoQueueListener) error {
	return this.BeegoQueueStoreInterface.Subscribe(topicId, this.WrapListener(listener))
}
