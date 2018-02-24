package queue

import (
	. "github.com/fishedee/app/log"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type RedisQueueStore struct {
	redisPool        *redis.Pool
	log              Log
	poolsize         int
	prefix           string
	savePath         string
	password         string
	dbNum            int
	retryInterval    int
	consumeListeners sync.Map
	isClose          int32
}

var MAX_POOL_SIZE = 100

func NewRedisQueue(log Log, config QueueStoreConfig) (QueueStoreInterface, error) {
	var savePath string
	var password string
	var dbNum int
	var poolsize int
	configs := strings.Split(config.SavePath, ",")
	if len(configs) > 0 {
		savePath = configs[0]
	}
	if len(configs) > 1 {
		poolsizeInner, err := strconv.Atoi(configs[1])
		if err != nil || poolsizeInner <= 0 {
			poolsize = MAX_POOL_SIZE
		} else {
			poolsize = poolsizeInner
		}
	} else {
		poolsize = MAX_POOL_SIZE
	}
	if len(configs) > 2 {
		password = configs[2]
	}
	if len(configs) > 3 {
		dbnumInt, err := strconv.Atoi(configs[3])
		if err != nil || dbnumInt < 0 {
			dbNum = 0
		} else {
			dbNum = dbnumInt
		}
	} else {
		dbNum = 0
	}
	if config.RetryInterval == 0 {
		config.RetryInterval = 5
	}

	result := &RedisQueueStore{
		prefix:        config.SavePrefix,
		savePath:      savePath,
		password:      password,
		dbNum:         dbNum,
		poolsize:      poolsize,
		retryInterval: config.RetryInterval,
		log:           log,
	}
	var err error
	result.redisPool, err = result.getConnectPool()
	if err != nil {
		return nil, err
	}
	return NewBasicQueue(result), nil
}

func (this *RedisQueueStore) getConnectPool() (*redis.Pool, error) {
	poollist := &redis.Pool{
		MaxIdle:     this.poolsize,
		IdleTimeout: 240 * time.Second,
		Dial:        this.getConnect,
	}
	return poollist, poollist.Get().Err()
}
func (this *RedisQueueStore) getConnect() (redis.Conn, error) {
	c, err := redis.DialTimeout(
		"tcp",
		this.savePath,
		time.Second,
		time.Second*12,
		time.Second,
	)
	if err != nil {
		return nil, err
	}
	if this.password != "" {
		if _, err := c.Do("AUTH", this.password); err != nil {
			c.Close()
			return nil, err
		}
	}
	_, err = c.Do("SELECT", this.dbNum)
	if err != nil {
		c.Close()
		return nil, err
	}
	return c, err
}
func (this *RedisQueueStore) Produce(topicId string, data []byte) error {
	c := this.redisPool.Get()
	defer c.Close()

	_, err := c.Do("LPUSH", this.prefix+topicId, data)
	if err != nil {
		return err
	}
	return nil
}

func (this *RedisQueueStore) consumeData(connect redis.Conn, topicId string, timeout int) ([]byte, error) {
	var topic interface{}
	var data interface{}

	reply, err := redis.Values(connect.Do("BRPOP", this.prefix+topicId, timeout))
	if err == redis.ErrNil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if reply == nil {
		return nil, nil
	}
	reply, err = redis.Scan(reply, &topic, &data)
	if err != nil {
		return nil, err
	}
	return data.([]byte), nil
}

func (this *RedisQueueStore) singleConsume(topicId string, listener QueueListener) error {
	conn, err := this.getConnect()
	if err != nil {
		return err
	}
	this.consumeListeners.Store(conn, true)
	defer func() {
		conn.Close()
		this.consumeListeners.Delete(conn)
	}()
	for {
		data, err := this.consumeData(conn, topicId, 5)
		if err != nil {
			return err
		}
		if data == nil {
			continue
		} else {
			listener(data)
		}
	}
}

func (this *RedisQueueStore) Consume(topicId string, listener QueueListener) error {
	go func() {
		for {
			err := this.singleConsume(topicId, listener)
			isExit := atomic.LoadInt32(&this.isClose)
			if isExit == 0 {
				this.log.Critical("Queue Redis consume error :%v, will be retry in %v seconds", err, this.retryInterval)
				time.Sleep(time.Duration(int(time.Second) * this.retryInterval))
			} else {
				return
			}
		}
	}()
	return nil
}

func (this *RedisQueueStore) Close() {
	atomic.StoreInt32(&this.isClose, 1)
	this.redisPool.Close()
	this.consumeListeners.Range(func(key, value interface{}) bool {
		conn := key.(redis.Conn)
		conn.Close()
		return true
	})
}
