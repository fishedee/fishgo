package queue

const (
	QUEUE_UNKNOWN = iota
	QUEUE_PUBLISH_SUBSCRIBE
	QUEUE_PRODUCE_CONSUME
)

type QueueListener func(argv []byte)

type QueueStoreInterface interface {
	Produce(topicId string, data []byte) error
	Consume(topicId string, listener QueueListener) error
	Publish(topicId string, data []byte) error
	Subscribe(topicId string, listener QueueListener) error
	Close()
}

type QueueStoreBasicInterface interface {
	Produce(topicId string, data []byte) error
	Consume(topicId string, listener QueueListener) error
	Close()
}

type QueueStoreConfig struct {
	SavePath      string
	SavePrefix    string
	RetryInterval int
}
