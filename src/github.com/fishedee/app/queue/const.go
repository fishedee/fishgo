package queue

type queueStoreListener func(argv []byte)

type queueStoreInterface interface {
	Produce(topicId string, data []byte) error
	Consume(topicId string, queue string, poolSize int, listener queueStoreListener) error
	Run() error
	Close()
}
