package queue

import (
	"errors"
	"fmt"
	"reflect"

	. "github.com/fishedee/app/log"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
)

type Queue interface {
	Produce(topicId string, data ...interface{}) error
	MustProduce(topicId string, data ...interface{})

	Consume(topicId string, queue string, poolSize int, listener interface{}) error
	MustConsume(topicId string, queue string, poolSize int, listener interface{})

	Run() error
	Close()
}

type QueueConfig struct {
	SavePath      string `config:"savepath"`
	SavePrefix    string `config:"saveprefix"`
	Driver        string `config:"driver"`
	Debug         bool   `config:"debug"`
	RetryInterval int    `config:"retryinterval"`
}

type queueImplement struct {
	store  queueStoreInterface
	log    Log
	config QueueConfig
}

func NewQueue(log Log, config QueueConfig) (Queue, error) {
	if config.Driver == "" {
		return nil, errors.New("invalid queue driver empty!")
	} else if config.Driver == "memory" {
		store, err := newMemoryQueue(log, config)
		if err != nil {
			return nil, err
		}
		return &queueImplement{
			log:    log,
			config: config,
			store:  store,
		}, nil
	} else if config.Driver == "redis" {
		store, err := newRedisQueue(log, config)
		if err != nil {
			return nil, err
		}
		return &queueImplement{
			log:    log,
			config: config,
			store:  store,
		}, nil
	} else {
		return nil, errors.New("invalid queue config " + config.Driver)
	}
}

func (this *queueImplement) encodeData(data []interface{}) ([]byte, error) {
	dataByte, err := EncodeJson(data)
	if err != nil {
		return nil, err
	}
	return dataByte, nil
}

func (this *queueImplement) decodeData(dataByte []byte, dataType []reflect.Type) ([]reflect.Value, error) {
	//读取数据
	result := []interface{}{}
	for _, singleDataType := range dataType {
		result = append(result, reflect.New(singleDataType).Interface())
	}
	err := DecodeJson(dataByte, &result)
	if err != nil {
		return nil, errors.New(err.Error() + "," + string(dataByte))
	}

	//构建参数
	valueResult := []reflect.Value{}
	for i := 0; i != len(dataType); i++ {
		if i >= len(result) {
			return nil, fmt.Errorf("call with %d argument function for %d argument, output: %v, input: %v", len(dataType), len(result), dataType, valueResult)
		}
		valueResult = append(valueResult, reflect.ValueOf(result[i]).Elem())
	}
	return valueResult, nil
}

func (this *queueImplement) WrapExceptionListener(listener interface{}, debugPrefix string) (queueStoreListener, error) {
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
	return func(data []byte) {
		if this.config.Debug {
			this.log.Debug("[Queue Consume] %v:%v", debugPrefix, string(data))
		}
		defer CatchCrash(func(exception Exception) {
			this.log.Critical("QueueTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
		})
		dataResult, err := this.decodeData(data, listenerInType)
		if err != nil {
			panic(err)
		}
		listenerValue.Call(dataResult)
	}, nil
}

func (this *queueImplement) Produce(topicId string, data ...interface{}) error {
	dataResult, err := this.encodeData(data)
	if err != nil {
		return err
	}
	err = this.store.Produce(topicId, dataResult)
	if err != nil {
		return err
	}
	if this.config.Debug {
		this.log.Debug("[Queue Produce] %v:%v", topicId, string(dataResult))
	}
	return nil
}

func (this *queueImplement) MustProduce(topicId string, data ...interface{}) {
	err := this.Produce(topicId, data...)
	if err != nil {
		panic(err)
	}
}

func (this *queueImplement) Consume(topicId string, queue string, poolSize int, listener interface{}) error {
	if poolSize <= 0 {
		poolSize = 1
	}
	storeListener, err := this.WrapExceptionListener(listener, topicId+"-"+queue)
	if err != nil {
		return err
	}
	return this.store.Consume(topicId, queue, poolSize, storeListener)
}

func (this *queueImplement) MustConsume(topicId string, queue string, poolSize int, listener interface{}) {
	err := this.Consume(topicId, queue, poolSize, listener)
	if err != nil {
		panic(err)
	}
}

func (this *queueImplement) Run() error {
	return this.store.Run()
}

func (this *queueImplement) Close() {
	this.store.Close()
}
