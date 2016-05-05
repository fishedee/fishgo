package web

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/fishedee/language"
	. "github.com/fishedee/web/util_queue"
	"reflect"
	"strconv"
)

type Queue interface {
	WithLog(log Log) Queue
	Produce(topicId string, data ...interface{}) error
	Consume(topicId string, listener interface{}) error
	ConsumeInPool(topicId string, listener interface{}, poolSize int) error
	Publish(topicId string, data ...interface{}) error
	Subscribe(topicId string, listener interface{}) error
	SubscribeInPool(topicId string, listener interface{}, poolSize int) error
}

type QueueConfig struct {
	SavePath   string
	SavePrefix string
	Driver     string
	PoolSize   int
}

type queueImplement struct {
	store    QueueStoreInterface
	Log      Log
	poolSize int
}

func NewQueue(config QueueConfig) (Queue, error) {
	if config.Driver == "" {
		return nil, nil
	} else if config.Driver == "memory" {
		queue, err := NewMemoryQueue(QueueStoreConfig{})
		if err != nil {
			return nil, err
		}
		return &queueImplement{
			store:    queue,
			poolSize: config.PoolSize,
		}, nil
	} else if config.Driver == "redis" {
		queue, err := NewRedisQueue(QueueStoreConfig{
			SavePath:   config.SavePath,
			SavePrefix: config.SavePrefix,
		})
		if err != nil {
			return nil, err
		}
		return &queueImplement{
			store:    queue,
			poolSize: config.PoolSize,
		}, nil
	} else {
		return nil, errors.New("invalid memory config " + config.Driver)
	}
}

func NewQueueFromConfig(configName string) (Queue, error) {
	driver := globalBasic.Config.GetString(configName + "driver")
	savepath := globalBasic.Config.GetString(configName + "savepath")
	saveprefix := globalBasic.Config.GetString(configName + "saveprefix")
	poolsizeStr := globalBasic.Config.GetString(configName + "poolsize")
	poolsize, _ := strconv.Atoi(poolsizeStr)

	queueConfig := QueueConfig{}
	queueConfig.Driver = driver
	queueConfig.SavePath = savepath
	queueConfig.SavePrefix = saveprefix
	queueConfig.PoolSize = poolsize
	return NewQueue(queueConfig)
}

func (this *queueImplement) WithLog(log Log) Queue {
	newQueueManager := *this
	newQueueManager.Log = log
	return &newQueueManager
}

func (this *queueImplement) EncodeData(data []interface{}) ([]byte, error) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return dataByte, nil
}

func (this *queueImplement) DecodeData(dataByte []byte, dataType []reflect.Type) ([]reflect.Value, error) {
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

func (this *queueImplement) WrapData(data []interface{}) (interface{}, error) {
	return this.EncodeData(data)
}

func (this *queueImplement) WrapPoolListener(listener QueueListener, poolSize int) QueueListener {
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

func (this *queueImplement) WrapExceptionListener(listener interface{}) (QueueListener, error) {
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

func (this *queueImplement) Produce(topicId string, data ...interface{}) error {
	dataResult, err := this.WrapData(data)
	if err != nil {
		return err
	}
	return this.store.Produce(topicId, dataResult)
}

func (this *queueImplement) Consume(topicId string, listener interface{}) error {
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

func (this *queueImplement) ConsumeInPool(topicId string, listener interface{}, poolSize int) error {
	listenerResult, err := this.WrapExceptionListener(listener)
	if err != nil {
		return err
	}
	if this.poolSize != 0 {
		poolSize = this.poolSize
	}
	return this.store.Consume(topicId, this.WrapPoolListener(listenerResult, poolSize))
}

func (this *queueImplement) Publish(topicId string, data ...interface{}) error {
	dataResult, err := this.WrapData(data)
	if err != nil {
		return err
	}
	return this.store.Publish(topicId, dataResult)
}

func (this *queueImplement) Subscribe(topicId string, listener interface{}) error {
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

func (this *queueImplement) SubscribeInPool(topicId string, listener interface{}, poolSize int) error {
	listenerResult, err := this.WrapExceptionListener(listener)
	if err != nil {
		return err
	}
	if this.poolSize != 0 {
		poolSize = this.poolSize
	}
	return this.store.Subscribe(topicId, this.WrapPoolListener(listenerResult, poolSize))
}
