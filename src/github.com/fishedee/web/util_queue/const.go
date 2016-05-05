package util_queue

const (
	QUEUE_UNKNOWN = iota
	QUEUE_PUBLISH_SUBSCRIBE
	QUEUE_PRODUCE_CONSUME
)

type QueueListener func(argv interface{}) error

type QueueStoreInterface interface {
	Produce(topicId string, data interface{}) error
	Consume(topicId string, listener QueueListener) error
	Publish(topicId string, data interface{}) error
	Subscribe(topicId string, listener QueueListener) error
}

type QueueStoreBasicInterface interface {
	Produce(topicId string, data interface{}) error
	Consume(topicId string, listener QueueListener) error
}

type QueueStoreConfig struct {
	SavePath   string
	SavePrefix string
}
