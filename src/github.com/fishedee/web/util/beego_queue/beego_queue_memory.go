package beego_queue

import (
	"errors"
	"sync"
)

type MemoryQueuePushPopStore struct {
	listener BeegoQueueListener
}

type MemoryQueueStore struct {
	mapPushPopStore map[string]MemoryQueuePushPopStore
	mutex           sync.Mutex
}

func NewMemoryQueue(BeegoQueueStoreConfig) (BeegoQueueStoreInterface, error) {
	result := &MemoryQueueStore{}
	result.mapPushPopStore = map[string]MemoryQueuePushPopStore{}
	return NewBasicQueue(result), nil
}

func (this *MemoryQueueStore) Produce(topicId string, data interface{}) error {
	this.mutex.Lock()
	result, ok := this.mapPushPopStore[topicId]
	this.mutex.Unlock()
	if !ok {
		return errors.New("empty listener")
	}
	result.listener(data)
	return nil
}

func (this *MemoryQueueStore) Consume(topicId string, listener BeegoQueueListener) error {
	this.mutex.Lock()
	this.mapPushPopStore[topicId] = MemoryQueuePushPopStore{
		listener: listener,
	}
	this.mutex.Unlock()
	return nil
}
