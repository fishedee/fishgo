package beego_queue

import (
	"errors"
	"sync"
)

type MemoryQueuePushPopStore struct {
	listener BeegoQueueListener
}

type MemoryQueuePubSubStore struct {
	listener []BeegoQueueListener
}

type MemoryQueueStore struct {
	mapPushPopStore map[string]MemoryQueuePushPopStore
	mapPubSubStore  map[string]MemoryQueuePubSubStore
	mutex           sync.Mutex
}

func NewMemoryQueue() (*MemoryQueueStore, error) {
	result := &MemoryQueueStore{}
	result.mapPushPopStore = map[string]MemoryQueuePushPopStore{}
	result.mapPubSubStore = map[string]MemoryQueuePubSubStore{}
	return result, nil
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

func (this *MemoryQueueStore) Publish(topicId string, data interface{}) error {
	this.mutex.Lock()
	result, ok := this.mapPubSubStore[topicId]
	this.mutex.Unlock()
	if !ok {
		return nil
	}
	for _, singleListener := range result.listener {
		singleListener(data)
	}
	return nil
}

func (this *MemoryQueueStore) Subscribe(topicId string, listener BeegoQueueListener) error {
	this.mutex.Lock()
	result, ok := this.mapPubSubStore[topicId]
	if !ok {
		result.listener = []BeegoQueueListener{}
	}
	result.listener = append(result.listener, listener)
	this.mapPubSubStore[topicId] = result
	this.mutex.Unlock()
	return nil
}
