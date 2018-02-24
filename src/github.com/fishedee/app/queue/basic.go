package queue

import (
	"sync"
)

type BasicAsyncQueuePubSubStore struct {
	listener []QueueListener
}

type BasicQueueStore struct {
	QueueStoreBasicInterface
	mapPubSubStore map[string]*BasicAsyncQueuePubSubStore
	mutex          sync.RWMutex
}

func NewBasicQueue(target QueueStoreBasicInterface) *BasicQueueStore {
	return &BasicQueueStore{
		QueueStoreBasicInterface: target,
		mapPubSubStore:           map[string]*BasicAsyncQueuePubSubStore{},
	}
}

func (this *BasicQueueStore) Publish(topicId string, data []byte) error {
	this.mutex.RLock()
	_, ok := this.mapPubSubStore[topicId]
	this.mutex.RUnlock()
	if !ok {
		return nil
	}
	return this.Produce(topicId, data)
}

func (this *BasicQueueStore) subscribeInner(topicId string, single *BasicAsyncQueuePubSubStore) error {
	return this.Consume(topicId, func(argv []byte, err error) {
		listeners := single.listener

		for _, singleListener := range listeners {
			singleListener(argv, err)
		}
	})
}

func (this *BasicQueueStore) Subscribe(topicId string, listener QueueListener) error {
	this.mutex.Lock()
	result, ok := this.mapPubSubStore[topicId]
	if !ok {
		result = &BasicAsyncQueuePubSubStore{}
		this.mapPubSubStore[topicId] = result
	}
	result.listener = append(result.listener, listener)
	this.mutex.Unlock()

	if !ok {
		//第一次订阅
		return this.subscribeInner(topicId, result)
	} else {
		//非第一次的订阅
		return nil
	}
}
