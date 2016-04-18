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
	"strconv"
)

type QueueManagerConfig struct {
	BeegoQueueStoreConfig
	Driver   string
	PoolSize int
}

type QueueManager struct {
	store    BeegoQueueStoreInterface
	Log      *LogManager
	poolSize int
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
			store:    queue,
			poolSize: config.PoolSize,
		}, nil
	} else if config.Driver == "redis" {
		queue, err := NewRedisQueue(config.BeegoQueueStoreConfig)
		if err != nil {
			return nil, err
		}
		return &QueueManager{
			store:    queue,
			poolSize: config.PoolSize,
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
	poolsizeStr := beego.AppConfig.String(configName + "poolsize")
	poolsize, _ := strconv.Atoi(poolsizeStr)

	queueConfig := QueueManagerConfig{}
	queueConfig.Driver = driver
	queueConfig.SavePath = savepath
	queueConfig.SavePrefix = saveprefix
	queueConfig.PoolSize = poolsize
	return NewQueueManager(queueConfig)
}

func NewQueueManagerFromConfig(configName string) (*QueueManager, error) {
	result, err := newQueueManagerFromConfigMemory.Call(configName)
	if err != nil {
		return nil, err
	}
	return result.(*QueueManager), err
}

func NewQueueManagerWithLog(log *LogManager, queue *QueueManager) *QueueManager {
	if queue == nil {
		return nil
	} else {
		return &QueueManager{
			store: queue.store,
			Log:   log,
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
		return nil, errors.New(err.Error() + "," + string(dataByte))
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

func (this *QueueManager) WrapPoolListener(listener BeegoQueueListener, poolSize int) BeegoQueueListener {
	if poolSize <= 0 {
		return func(data interface{}) (lastError error) {
			go listener(data)
			return nil
		}
	} else if poolSize == 1 {
		return listener
	} else {
		chanConsume := make(chan bool, poolSize)
		for i := 0; i != poolSize; i++ {
			chanConsume <- true
		}
		return func(data interface{}) (lastError error) {
			<-chanConsume
			go listener(data)
			chanConsume <- true
			return nil
		}
	}
}

func (this *QueueManager) WrapExceptionListener(listener interface{}) (BeegoQueueListener, error) {
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
			lastError = exception
		})
		defer Catch(func(exception Exception) {
			this.Log.Error("QueueTask Error Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
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
	listenerResult, err := this.WrapExceptionListener(listener)
	if err != nil {
		return err
	}
	poolSize := 0
	if this.poolSize != 0 {
		poolSize = this.poolSize
	}
	return this.store.Consume(topicId, this.WrapPoolListener(listenerResult, poolSize))
}

func (this *QueueManager) ConsumeInPool(topicId string, listener interface{}, poolSize int) error {
	listenerResult, err := this.WrapExceptionListener(listener)
	if err != nil {
		return err
	}
	if this.poolSize != 0 {
		poolSize = this.poolSize
	}
	return this.store.Consume(topicId, this.WrapPoolListener(listenerResult, poolSize))
}

func (this *QueueManager) Publish(topicId string, data ...interface{}) error {
	dataResult, err := this.WrapData(data)
	if err != nil {
		return err
	}
	return this.store.Publish(topicId, dataResult)
}

func (this *QueueManager) Subscribe(topicId string, listener interface{}) error {
	listenerResult, err := this.WrapExceptionListener(listener)
	if err != nil {
		return err
	}
	poolSize := 0
	if this.poolSize != 0 {
		poolSize = this.poolSize
	}
	return this.store.Subscribe(topicId, this.WrapPoolListener(listenerResult, poolSize))
}

func (this *QueueManager) SubscribeInPool(topicId string, listener interface{}, poolSize int) error {
	listenerResult, err := this.WrapExceptionListener(listener)
	if err != nil {
		return err
	}
	if this.poolSize != 0 {
		poolSize = this.poolSize
	}
	return this.store.Subscribe(topicId, this.WrapPoolListener(listenerResult, poolSize))
}
