package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
	. "github.com/fishedee/web/util/beego_queue"
	"reflect"
)

type QueueManagerConfig struct {
	BeegoQueueStoreConfig
	Driver string
}

type QueueManager struct {
	store   BeegoQueueStoreInterface
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
			store: queue,
		}, nil
	} else if config.Driver == "memory_async" {
		queue, err := NewMemoryAsyncQueue(config.BeegoQueueStoreConfig)
		if err != nil {
			return nil, err
		}
		return &QueueManager{
			store: queue,
		}, nil
	} else if config.Driver == "redis" {
		queue, err := NewRedisQueue(config.BeegoQueueStoreConfig)
		if err != nil {
			return nil, err
		}
		return &QueueManager{
			store: queue,
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
			store:   queue.store,
			Log:     log,
			Monitor: monitor,
		}
	}
}

func (this *QueueManager) EncodeData(data []interface{}) ([]byte, error) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return dataByte, nil
}

func (this *QueueManager) DecodeData(dataByte []byte, dataType []reflect.Type) ([]reflect.Value, error) {
	result := []interface{}{}
	for _, singleDataType := range dataType {
		result = append(result, reflect.New(singleDataType).Interface())
	}
	err := json.Unmarshal(dataByte, &result)
	if err != nil {
		return nil, err
	}
	valueResult := []reflect.Value{}
	for i := 0; i != len(dataType); i++ {
		if i >= len(result) {
			panic(fmt.Sprintf("call with %d argument function for %d argument", len(dataType), len(result)))
		}
		valueResult = append(valueResult, reflect.ValueOf(result[i]).Elem())
	}
	return valueResult, nil
}

func (this *QueueManager) WrapData(data []interface{}) (interface{}, error) {
	return this.EncodeData(data)
}

func (this *QueueManager) WrapListener(listener interface{}) (BeegoQueueListener, error) {
	listenerType := reflect.TypeOf(listener)
	listenerValue := reflect.ValueOf(listener)
	if listenerType.Kind() != reflect.Func {
		return nil, errors.New("listener type is not a function")
	}
	listenerInType := []reflect.Type{}
	for i := 0; i != listenerType.NumIn(); i++ {
		listenerInType = append(
			listenerInType,
			listenerType.In(i),
		)
	}
	return func(data interface{}) (lastError error) {
		defer CatchCrash(func(exception Exception) {
			this.Log.Critical("QueueTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
			if this.Monitor != nil {
				this.Monitor.AscCriticalCount()
			}
			lastError = exception
		})
		defer Catch(func(exception Exception) {
			this.Log.Error("QueueTask Error Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
			if this.Monitor != nil {
				this.Monitor.AscErrorCount()
			}
			lastError = exception
		})
		argvError, ok := data.(error)
		if ok {
			panic(argvError)
		}
		dataResult, err := this.DecodeData(data.([]byte), listenerInType)
		if err != nil {
			Throw(1, err.Error())
		}
		listenerValue.Call(dataResult)
		return nil
	}, nil
}

func (this *QueueManager) Produce(topicId string, data ...interface{}) error {
	dataResult, err := this.WrapData(data)
	if err != nil {
		return err
	}
	return this.store.Produce(topicId, dataResult)
}

func (this *QueueManager) Consume(topicId string, listener interface{}) error {
	listenerResult, err := this.WrapListener(listener)
	if err != nil {
		return err
	}
	return this.store.Consume(topicId, listenerResult)
}

func (this *QueueManager) Publish(topicId string, data ...interface{}) error {
	dataResult, err := this.WrapData(data)
	if err != nil {
		return err
	}
	return this.store.Publish(topicId, dataResult)
}

func (this *QueueManager) Subscribe(topicId string, listener interface{}) error {
	listenerResult, err := this.WrapListener(listener)
	if err != nil {
		return err
	}
	return this.store.Subscribe(topicId, listenerResult)
}
