package web

import (
	"errors"
	"fmt"
	. "github.com/fishedee/encoding"
	. "github.com/fishedee/language"
	. "github.com/fishedee/util"
	. "github.com/fishedee/web/util_queue"
	"reflect"
	"strconv"
)

type Queue interface {
	WithLogAndContext(log Log, ctx Context) Queue
	Produce(topicId string, data ...interface{})
	Consume(topicId string, listener interface{})
	ConsumeInPool(topicId string, listener interface{}, poolSize int)
	Publish(topicId string, data ...interface{})
	Subscribe(topicId string, listener interface{})
	SubscribeInPool(topicId string, listener interface{}, poolSize int)
	Close()
}

type QueueConfig struct {
	SavePath   string
	SavePrefix string
	Driver     string
	PoolSize   int
	Debug      bool
}

type queueImplement struct {
	store     QueueStoreInterface
	Log       Log
	Ctx       Context
	poolSize  int
	debug     bool
	closeFunc *CloseFunc
}

func NewQueue(config QueueConfig) (Queue, error) {
	if config.Driver == "" {
		return nil, nil
	} else if config.Driver == "memory" {
		closeFunc := NewCloseFunc()
		queue, err := NewMemoryQueue(closeFunc, QueueStoreConfig{})
		if err != nil {
			return nil, err
		}
		return &queueImplement{
			store:     queue,
			poolSize:  config.PoolSize,
			debug:     config.Debug,
			closeFunc: closeFunc,
		}, nil
	} else if config.Driver == "redis" {
		closeFunc := NewCloseFunc()
		queue, err := NewRedisQueue(closeFunc, QueueStoreConfig{
			SavePath:   config.SavePath,
			SavePrefix: config.SavePrefix,
		})
		if err != nil {
			return nil, err
		}
		return &queueImplement{
			store:     queue,
			poolSize:  config.PoolSize,
			debug:     config.Debug,
			closeFunc: closeFunc,
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
	debugStr := globalBasic.Config.GetString(configName + "debug")
	debug, _ := strconv.ParseBool(debugStr)

	queueConfig := QueueConfig{}
	queueConfig.Driver = driver
	queueConfig.SavePath = savepath
	queueConfig.SavePrefix = saveprefix
	queueConfig.PoolSize = poolsize
	queueConfig.Debug = debug
	return NewQueue(queueConfig)
}

func (this *queueImplement) WithLogAndContext(log Log, ctx Context) Queue {
	newQueueManager := *this
	newQueueManager.Log = log
	newQueueManager.Ctx = ctx
	return &newQueueManager
}

func (this *queueImplement) EncodeData(data []interface{}) ([]byte, error) {
	ctxRequest, err := this.Ctx.SerializeRequest()
	if err != nil {
		return nil, err
	}
	data = append([]interface{}{ctxRequest}, data...)
	dataByte, err := EncodeJson(data)
	if err != nil {
		return nil, err
	}
	return dataByte, nil
}

func (this *queueImplement) DecodeData(dataByte []byte, dataType []reflect.Type) ([]reflect.Value, error) {
	//读取数据
	result := []interface{}{}
	var ctxResult ContextSerializeRequest
	for singleDataIndex, singleDataType := range dataType {
		if singleDataIndex == 0 {
			result = append(result, &ctxResult)
		} else {
			result = append(result, reflect.New(singleDataType).Interface())
		}

	}
	err := DecodeJson(dataByte, &result)
	if err != nil {
		return nil, errors.New(err.Error() + "," + string(dataByte))
	}

	//构建参数
	basic := initEmptyBasic(nil)
	target := reflect.New(dataType[0].Elem())
	err = basic.Ctx.DeSerializeRequest(ctxResult)
	if err != nil {
		return nil, err
	}
	injectIoc(target, basic)
	valueResult := []reflect.Value{}
	valueResult = append(valueResult, target)
	for i := 1; i != len(dataType); i++ {
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
			this.closeFunc.IncrCloseCounter()
			go func() {
				defer this.closeFunc.DecrCloseCounter()
				listener(data)
			}()
			return nil
		}
	} else if poolSize == 1 {
		return func(data interface{}) error {
			this.closeFunc.IncrCloseCounter()
			defer this.closeFunc.DecrCloseCounter()
			return listener(data)
		}
	} else {
		chanConsume := make(chan bool, poolSize)
		for i := 0; i != poolSize; i++ {
			chanConsume <- true
		}
		return func(data interface{}) (lastError error) {
			this.closeFunc.IncrCloseCounter()
			<-chanConsume
			go func() {
				defer func() {
					chanConsume <- true
				}()
				defer this.closeFunc.DecrCloseCounter()
				listener(data)
			}()
			return nil
		}
	}
}

func (this *queueImplement) WrapExceptionListener(listener interface{}, topicId string, useplace string) (QueueListener, error) {
	listenerType := reflect.TypeOf(listener)
	listenerValue := reflect.ValueOf(listener)
	if listenerType.Kind() != reflect.Func {
		return nil, errors.New("listener type is not a function")
	}
	if listenerType.NumIn() == 0 {
		return nil, errors.New("listener should has at last a argument")
	}
	if listenerType.In(0).Kind() != reflect.Ptr {
		return nil, errors.New("listener first argument is not a ptr")
	}
	listenerInType := []reflect.Type{}
	for i := 0; i != listenerType.NumIn(); i++ {
		listenerInType = append(
			listenerInType,
			listenerType.In(i),
		)
	}
	return func(data interface{}) (lastError error) {
		if this.debug {
			this.Log.Debug("[Queue %v] %v:%v", useplace, topicId, string(data.([]byte)))
		}
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

func (this *queueImplement) Produce(topicId string, data ...interface{}) {
	defer CatchCrash(func(exception Exception) {
		this.Log.Critical("QueueTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	dataResult, err := this.WrapData(data)
	if err != nil {
		panic(err)
	}
	err = this.store.Produce(topicId, dataResult)
	if err != nil {
		panic(err)
	}
	if this.debug {
		this.Log.Debug("[Queue Produce] %v:%v", topicId, string(dataResult.([]byte)))
	}
}

func (this *queueImplement) Consume(topicId string, listener interface{}) {
	defer CatchCrash(func(exception Exception) {
		this.Log.Critical("QueueTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	listenerResult, err := this.WrapExceptionListener(listener, topicId, "Consume")
	if err != nil {
		panic(err)
	}
	poolSize := 0
	if this.poolSize != 0 {
		poolSize = this.poolSize
	}
	err = this.store.Consume(topicId, this.WrapPoolListener(listenerResult, poolSize))
	if err != nil {
		panic(err)
	}
}

func (this *queueImplement) ConsumeInPool(topicId string, listener interface{}, poolSize int) {
	defer CatchCrash(func(exception Exception) {
		this.Log.Critical("QueueTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	listenerResult, err := this.WrapExceptionListener(listener, topicId, "ConsumeInPool")
	if err != nil {
		panic(err)
	}
	if this.poolSize != 0 {
		poolSize = this.poolSize
	}
	err = this.store.Consume(topicId, this.WrapPoolListener(listenerResult, poolSize))
	if err != nil {
		panic(err)
	}
}

func (this *queueImplement) Publish(topicId string, data ...interface{}) {
	defer CatchCrash(func(exception Exception) {
		this.Log.Critical("QueueTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	dataResult, err := this.WrapData(data)
	if err != nil {
		panic(err)
	}
	err = this.store.Publish(topicId, dataResult)
	if err != nil {
		panic(err)
	}
	if this.debug {
		this.Log.Debug("[Queue Publish] %v:%v", topicId, string(dataResult.([]byte)))
	}
}

func (this *queueImplement) Subscribe(topicId string, listener interface{}) {
	defer CatchCrash(func(exception Exception) {
		this.Log.Critical("QueueTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	listenerResult, err := this.WrapExceptionListener(listener, topicId, "Subscribe")
	if err != nil {
		panic(err)
	}
	poolSize := 0
	if this.poolSize != 0 {
		poolSize = this.poolSize
	}
	err = this.store.Subscribe(topicId, this.WrapPoolListener(listenerResult, poolSize))
	if err != nil {
		panic(err)
	}
}

func (this *queueImplement) SubscribeInPool(topicId string, listener interface{}, poolSize int) {
	defer CatchCrash(func(exception Exception) {
		this.Log.Critical("QueueTask Crash Code:[%d] Message:[%s]\nStackTrace:[%s]", exception.GetCode(), exception.GetMessage(), exception.GetStackTrace())
	})
	listenerResult, err := this.WrapExceptionListener(listener, topicId, "SubscribeInPool")
	if err != nil {
		panic(err)
	}
	if this.poolSize != 0 {
		poolSize = this.poolSize
	}
	err = this.store.Subscribe(topicId, this.WrapPoolListener(listenerResult, poolSize))
	if err != nil {
		panic(err)
	}
}

func (this *queueImplement) Close() {
	this.closeFunc.Close()
}
