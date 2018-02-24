package queue

const (
	QUEUE_UNKNOWN = iota
	QUEUE_PUBLISH_SUBSCRIBE
	QUEUE_PRODUCE_CONSUME
)

type QueueListener func(argv []byte, err error)

type queueInnerListener func(argv []byte)

type QueueStoreInterface interface {
	Produce(topicId string, data []byte) error
	Consume(topicId string, listener QueueListener) error
	Publish(topicId string, data []byte) error
	Subscribe(topicId string, listener QueueListener) error
}

type QueueStoreBasicInterface interface {
	Produce(topicId string, data []byte) error
	Consume(topicId string, listener QueueListener) error
}

type QueueStoreConfig struct {
	SavePath   string
	SavePrefix string
}
