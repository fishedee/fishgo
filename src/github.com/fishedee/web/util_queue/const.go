package beego_queue

const (
	BEEGO_QUEUE_UNKNOWN = iota
	BEEGO_QUEUE_PUBLISH_SUBSCRIBE
	BEEGO_QUEUE_PRODUCE_CONSUME
)

type BeegoQueueListener func(argv interface{}) error

type BeegoQueueStoreInterface interface {
	Produce(topicId string, data interface{}) error
	Consume(topicId string, listener BeegoQueueListener) error
	Publish(topicId string, data interface{}) error
	Subscribe(topicId string, listener BeegoQueueListener) error
}

type BeegoQueueStoreBasicInterface interface {
	Produce(topicId string, data interface{}) error
	Consume(topicId string, listener BeegoQueueListener) error
}

type BeegoQueueStoreConfig struct {
	SavePath   string
	SavePrefix string
}
