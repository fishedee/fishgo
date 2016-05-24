package util_queue

import (
	. "github.com/fishedee/util"
	"sync"
)

type MemoryQueuePushPopStore struct {
	listener QueueListener
}

type MemoryQueueStore struct {
	mapPushPopStore map[string]MemoryQueuePushPopStore
	mutex           sync.Mutex
}

func NewMemoryQueue(closeFunc *CloseFunc, config QueueStoreConfig) (QueueStoreInterface, error) {
	result := &MemoryQueueStore{}
	result.mapPushPopStore = map[string]MemoryQueuePushPopStore{}
	return NewBasicQueue(result), nil
}

func (this *MemoryQueueStore) Produce(topicId string, data interface{}) error {
	this.mutex.Lock()
	result, ok := this.mapPushPopStore[topicId]
	this.mutex.Unlock()
	if !ok {
		return nil
	}
	result.listener(data)
	return nil
}

func (this *MemoryQueueStore) Consume(topicId string, listener QueueListener) error {
	this.mutex.Lock()
	this.mapPushPopStore[topicId] = MemoryQueuePushPopStore{
		listener: listener,
	}
	this.mutex.Unlock()
	return nil
}
