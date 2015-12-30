package util

import (
	"errors"
	"github.com/astaxie/beego"
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
	. "github.com/fishedee/web/util/beego_queue"
)

type QueueManagerConfig struct {
	BeegoQueueStoreConfig
	Driver string
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
		queue, err := NewMemoryQueue(config.BeegoQueueStoreConfig)
		if err != nil {
			return nil, err
		}
		return &QueueManager{
			BeegoQueueStoreInterface: queue,
		}, nil
	} else if config.Driver == "memory_async" {
		queue, err := NewMemoryAsyncQueue(config.BeegoQueueStoreConfig)
		if err != nil {
			return nil, err
		}
		return &QueueManager{
			BeegoQueueStoreInterface: queue,
		}, nil
	} else if config.Driver == "redis" {
		queue, err := NewRedisQueue(config.BeegoQueueStoreConfig)
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
	savepath := beego.AppConfig.String(configName + "savepath")
	saveprefix := beego.AppConfig.String(configName + "saveprefix")

	queueConfig := QueueManagerConfig{}
	queueConfig.Driver = driver
	queueConfig.SavePath = savepath
	queueConfig.SavePrefix = saveprefix
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
	return func(data interface{}) {
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
		argvError, ok := data.(error)
		if ok {
			panic(argvError)
		}
		listener(data)
	}
}

func (this *QueueManager) Consume(topicId string, listener BeegoQueueListener) error {
	return this.BeegoQueueStoreInterface.Consume(topicId, this.WrapListener(listener))
}

func (this *QueueManager) Subscribe(topicId string, listener BeegoQueueListener) error {
	return this.BeegoQueueStoreInterface.Subscribe(topicId, this.WrapListener(listener))
}
