package queue

import (
	. "github.com/fishedee/util"
	"sync"
)

type MemoryQueuePushPopStore struct {
	channel *InfiniteChannel
}

type MemoryQueueStore struct {
	mapPushPopStore map[string]MemoryQueuePushPopStore
	mutex           sync.RWMutex
	waitgroup       *sync.WaitGroup
}

func NewMemoryQueue(config QueueStoreConfig) (QueueStoreInterface, error) {
	result := &MemoryQueueStore{}
	result.mapPushPopStore = map[string]MemoryQueuePushPopStore{}
	result.waitgroup = &sync.WaitGroup{}
	return NewBasicQueue(result), nil
}

func (this *MemoryQueueStore) getTopicInfo(topicId string) MemoryQueuePushPopStore {
	this.mutex.RLock()
	result, ok := this.mapPushPopStore[topicId]
	this.mutex.RUnlock()
	if !ok {
		result = MemoryQueuePushPopStore{
			channel: NewInfiniteChannel(),
		}
		this.mutex.Lock()
		this.mapPushPopStore[topicId] = result
		this.mutex.Unlock()
	}
	return result
}
func (this *MemoryQueueStore) Produce(topicId string, data []byte) error {
	result := this.getTopicInfo(topicId)
	result.channel.Write(data)
	return nil
}

func (this *MemoryQueueStore) Consume(topicId string, listener QueueListener) error {
	result := this.getTopicInfo(topicId)
	this.waitgroup.Add(1)
	go func() {
		defer this.waitgroup.Done()
		for single := range result.channel.Read() {
			listener(single.([]byte))
		}
	}()
	return nil
}

func (this *MemoryQueueStore) Close() {
	this.mutex.RLock()
	for _, single := range this.mapPushPopStore {
		single.channel.Close()
	}
	this.mutex.RUnlock()

	this.waitgroup.Wait()
}
