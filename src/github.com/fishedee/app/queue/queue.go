package queue

import (
	"errors"
	"fmt"
	"reflect"

	. "github.com/fishedee/app/log"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
)

type Queue interface {
	Produce(topicId string, data ...interface{}) error
	MustProduce(topicId string, data ...interface{})

	Consume(topicId string, listener interface{}) error
	MustConsume(topicId string, listener interface{})

	ConsumeInPool(topicId string, listener interface{}, poolSize int) error
	MustConsumeInPool(topicId string, listener interface{}, poolSize int)

	Publish(topicId string, data ...interface{}) error
	MustPublish(topicId string, data ...interface{})

	Subscribe(topicId string, listener interface{}) error
	MustSubscribe(topicId string, listener interface{})

	SubscribeInPool(topicId string, listener interface{}, poolSize int) error
	MustSubscribeInPool(topicId string, listener interface{}, poolSize int)

	Run() error
	Close()
}

type QueueConfig struct {
	SavePath      string `config:"savepath"`
	SavePrefix    string `config:"saveprefix"`
	Driver        string `config:"driver"`
	PoolSize      int    `config:"poolsize"`
	Debug         bool   `config:"debug"`
	RetryInterval int    `config:"retryinterval"`
}

type queueImplement struct {
	store     QueueStoreInterface
	log       Log
	poolSize  int
	debug     bool
	closeFunc *CloseFunc
	exitChan  chan bool
}

func NewQueue(log Log, config QueueConfig) (Queue, error) {
	if config.Driver == "" {
		return nil, nil
	} else if config.Driver == "memory" {
		closeFunc := NewCloseFunc()
		queue, err := NewMemoryQueue(QueueStoreConfig{})
		if err != nil {
			return nil, err
		}
		return &queueImplement{
			log:       log,
			store:     queue,
			poolSize:  config.PoolSize,
			debug:     config.Debug,
			closeFunc: closeFunc,
			exitChan:  make(chan bool, 8),
		}, nil
	} else if config.Driver == "redis" {
		closeFunc := NewCloseFunc()
		queue, err := NewRedisQueue(log, QueueStoreConfig{
			SavePath:      config.SavePath,
			SavePrefix:    config.SavePrefix,
			RetryInterval: config.RetryInterval,
		})
		if err != nil {
			return nil, err
		}
		return &queueImplement{
			log:       log,
			store:     queue,
			poolSize:  config.PoolSize,
			debug:     config.Debug,
			closeFunc: closeFunc,
			exitChan:  make(chan bool, 8),
		}, nil
	} else {
		return nil, errors.New("invalid memory config " + config.Driver)
	}
}

func (this *queueImplement) EncodeData(data []interface{}) ([]byte, error) {
	dataByte, err := EncodeJson(data)
	if err != nil {
		return nil, err
	}
	return dataByte, nil
}

func (this *queueImplement) DecodeData(dataByte []byte, dataType []reflect.Type) ([]reflect.Value, error) {
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

func (this *queueImplement) WrapPoolListener(listener QueueListener, poolSize int) QueueListener {
	if poolSize <= 0 {
		return func(data []byte) {
			this.closeFunc.IncrCloseCounter()
			go func() {
				defer this.closeFunc.DecrCloseCounter()
				listener(data)
			}()
		}
	} else if poolSize == 1 {
		return func(data []byte) {
			this.closeFunc.IncrCloseCounter()
			defer this.closeFunc.DecrCloseCounter()
			listener(data)
		}
	} else {
		chanConsume := make(chan bool, poolSize)
		for i := 0; i != poolSize; i++ {
			chanConsume <- true
		}
		return func(data []byte) {
			this.closeFunc.IncrCloseCounter()
			<-chanConsume
			go func() {
				defer func() {
					chanConsume <- true
				}()
				defer this.closeFunc.DecrCloseCounter()
				listener(data)
			}()
		}
	}
}

func (this *queueImplement) WrapExceptionListener(listener interface{}, topicId string, debugPrefix string) (QueueListener, error) {
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
		if this.debug {
			this.log.Debug("[Queue %v] %v:%v", debugPrefix, topicId, string(data))
		}
		defer CatchCrash(func(exception Exception) {
			this.log.Critical("QueueTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
		})
		dataResult, err := this.DecodeData(data, listenerInType)
		if err != nil {
			panic(err)
		}
		listenerValue.Call(dataResult)
	}, nil
}

func (this *queueImplement) Produce(topicId string, data ...interface{}) error {
	dataResult, err := this.EncodeData(data)
	if err != nil {
		return err
	}
	err = this.store.Produce(topicId, dataResult)
	if err != nil {
		return err
	}
	if this.debug {
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

func (this *queueImplement) Consume(topicId string, listener interface{}) error {
	return this.ConsumeInPool(topicId, listener, this.poolSize)
}

func (this *queueImplement) MustConsume(topicId string, listener interface{}) {
	err := this.Consume(topicId, listener)
	if err != nil {
		panic(err)
	}
}

func (this *queueImplement) ConsumeInPool(topicId string, listener interface{}, poolSize int) error {
	listenerResult, err := this.WrapExceptionListener(listener, topicId, "Consume")
	if err != nil {
		return err
	}
	err = this.store.Consume(topicId, this.WrapPoolListener(listenerResult, poolSize))
	if err != nil {
		return err
	}
	return nil
}

func (this *queueImplement) MustConsumeInPool(topicId string, listener interface{}, poolSize int) {
	err := this.ConsumeInPool(topicId, listener, poolSize)
	if err != nil {
		panic(err)
	}
}

func (this *queueImplement) Publish(topicId string, data ...interface{}) error {
	dataResult, err := this.EncodeData(data)
	if err != nil {
		return err
	}
	err = this.store.Publish(topicId, dataResult)
	if err != nil {
		return err
	}
	if this.debug {
		this.log.Debug("[Queue Publish] %v:%v", topicId, string(dataResult))
	}
	return nil
}

func (this *queueImplement) MustPublish(topicId string, data ...interface{}) {
	err := this.Publish(topicId, data...)
	if err != nil {
		panic(err)
	}
}

func (this *queueImplement) Subscribe(topicId string, listener interface{}) error {
	return this.SubscribeInPool(topicId, listener, this.poolSize)
}

func (this *queueImplement) MustSubscribe(topicId string, listener interface{}) {
	err := this.Subscribe(topicId, listener)
	if err != nil {
		panic(err)
	}
}

func (this *queueImplement) SubscribeInPool(topicId string, listener interface{}, poolSize int) error {
	listenerResult, err := this.WrapExceptionListener(listener, topicId, "Subscribe")
	if err != nil {
		return err
	}
	err = this.store.Subscribe(topicId, this.WrapPoolListener(listenerResult, poolSize))
	if err != nil {
		return err
	}
	return nil
}

func (this *queueImplement) MustSubscribeInPool(topicId string, listener interface{}, poolSize int) {
	err := this.SubscribeInPool(topicId, listener, poolSize)
	if err != nil {
		panic(err)
	}
}

func (this *queueImplement) Run() error {
	<-this.exitChan
	this.store.Close()
	this.closeFunc.Close()
	return nil
}

func (this *queueImplement) Close() {
	this.exitChan <- true
}
