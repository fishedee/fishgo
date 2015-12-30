package beego_queue

import (
	"sync"
)

type MemoryAsyncQueuePushPopStore struct {
	channel dataChannel
}

type MemoryAsyncQueuePubSubStore struct {
	channels []dataChannel
}

type dataChannel chan interface{}

type MemoryAsyncQueueStore struct {
	mapPushPopStore map[string]MemoryAsyncQueuePushPopStore
	mapPubSubStore  map[string]MemoryAsyncQueuePubSubStore
	mutex           sync.Mutex
	chanSize        int
}

func NewMemoryAsyncQueue() (*MemoryAsyncQueueStore, error) {
	result := &MemoryAsyncQueueStore{}
	result.mapPushPopStore = map[string]MemoryAsyncQueuePushPopStore{}
	result.mapPubSubStore = map[string]MemoryAsyncQueuePubSubStore{}
	result.chanSize = 1024
	return result, nil
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

func (this *MemoryAsyncQueueStore) Publish(topicId string, data interface{}) error {
	this.mutex.Lock()
	result, ok := this.mapPubSubStore[topicId]
	this.mutex.Unlock()
	if !ok {
		return nil
	}
	for _, singleChannel := range result.channels {
		singleChannel <- data
	}
	return nil
}

func (this *MemoryAsyncQueueStore) Subscribe(topicId string, listener BeegoQueueListener) error {
	this.mutex.Lock()
	result, ok := this.mapPubSubStore[topicId]
	if !ok {
		result.channels = []dataChannel{}
	}
	newChannel := make(dataChannel, this.chanSize)
	result.channels = append(result.channels, newChannel)
	this.mapPubSubStore[topicId] = result
	this.mutex.Unlock()

	go func() {
		for {
			data := <-newChannel
			listener(data)
		}
	}()
	return nil
}
