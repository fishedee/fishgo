package beego_queue

import (
	"sync"
)

type MemoryAsyncQueuePushPopStore struct {
	channel dataChannel
}

type dataChannel chan interface{}

type MemoryAsyncQueueStore struct {
	mapPushPopStore map[string]MemoryAsyncQueuePushPopStore
	mutex           sync.Mutex
	chanSize        int
}

func NewMemoryAsyncQueue(BeegoQueueStoreConfig) (BeegoQueueStoreInterface, error) {
	result := &MemoryAsyncQueueStore{}
	result.mapPushPopStore = map[string]MemoryAsyncQueuePushPopStore{}
	result.chanSize = 1024
	return NewBasicAsyncQueue(result), nil
}

func (this *MemoryAsyncQueueStore) Produce(topicId string, data interface{}) error {
	this.mutex.Lock()
	result, ok := this.mapPushPopStore[topicId]
	if !ok {
		result.channel = make(dataChannel, this.chanSize)
		this.mapPushPopStore[topicId] = result
	}
	this.mutex.Unlock()

	result.channel <- data
	return nil
}

func (this *MemoryAsyncQueueStore) Consume(topicId string, listener BeegoQueueListener) error {
	this.mutex.Lock()
	result, ok := this.mapPushPopStore[topicId]
	if !ok {
		result.channel = make(dataChannel, this.chanSize)
		this.mapPushPopStore[topicId] = result
	}
	this.mutex.Unlock()

	go func() {
		for {
			data := <-result.channel
			listener(data)
		}
	}()
	return nil
}
