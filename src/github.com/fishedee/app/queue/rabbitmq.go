package queue

import (
	"errors"
	. "github.com/fishedee/app/log"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

type rabbitmqQueueStore struct {
	log       Log
	config    QueueConfig
	waitgroup *sync.WaitGroup
	closeChan chan bool
	exitChan  chan bool
	listener  sync.Map
}

func newRabbitmqQueueStore(log Log, config QueueConfig) (*rabbitmqQueueStore, error) {
	if config.RetryInterval == 0 {
		config.RetryInterval = 5
	}
	return &rabbitmqQueueStore{
		log:       log,
		config:    config,
		waitgroup: &sync.WaitGroup{},
		closeChan: make(chan bool, 16),
		exitChan:  make(chan bool, 16),
	}, nil
}

func (this *rabbitmqQueueStore) Produce(topicId string, data []byte) error {
	conn, err := amqp.Dial(this.config.SavePath)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Publish(
		topicId, // exchange
		"",      // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         data,
			DeliveryMode: 2,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

func (this *rabbitmqQueueStore) singleConsume(queue string, listener queueStoreListener) error {
	conn, err := amqp.Dial(this.config.SavePath)
	if err != nil {
		return err
	}
	defer func() {
		conn.Close()
		this.listener.Delete(conn)
	}()

	this.listener.Store(conn, true)

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Qos(
		1,    //prefetchCount
		0,    //prefetchSize
		true, //global
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	for {
		select {
		case d, isOk := <-msgs:
			if isOk == false {
				return errors.New("close consume")
			}
			listener(d.Body)
			err := d.Ack(false)
			if err != nil {
				return err
			}
		case _, _ = <-this.closeChan:
			return nil
		}
	}
}

func (this *rabbitmqQueueStore) buildTopic(topicId string, queue string) error {
	conn, err := amqp.Dial(this.config.SavePath)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		topicId,  // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		queue,   // queue name
		"",      // routing key
		topicId, // exchange
		false,   //no-wait
		nil)
	if err != nil {
		return err
	}

	return nil
}

func (this *rabbitmqQueueStore) Consume(topicId string, queue string, poolSize int, listener queueStoreListener) error {
	err := this.buildTopic(topicId, queue)
	if err != nil {
		return err
	}
	for i := 0; i < poolSize; i++ {
		this.waitgroup.Add(1)
		go func() {
			defer this.waitgroup.Done()
			for {
				err := this.singleConsume(queue, listener)
				if err != nil {
					select {
					case _, _ = <-this.closeChan:
						return
					default:
					}
					this.log.Critical("Queue Rabbitmq consume error :%v, will be retry in %v seconds", err, this.config.RetryInterval)
					sleepTime := int(time.Second) * this.config.RetryInterval
					timer := time.NewTimer(time.Duration(sleepTime))
					select {
					case _, _ = <-this.closeChan:
						return
					case _ = <-timer.C:
						break
					}
				} else {
					return
				}
			}
		}()
	}
	return nil
}

func (this *rabbitmqQueueStore) Run() error {
	_, _ = <-this.closeChan
	this.waitgroup.Wait()
	this.exitChan <- true
	return nil
}

func (this *rabbitmqQueueStore) Close() {
	close(this.closeChan)
	<-this.exitChan
}
