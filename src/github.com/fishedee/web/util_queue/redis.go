package util_queue

import (
	. "github.com/fishedee/util"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
	"time"
)

type RedisQueueStore struct {
	redisPool *redis.Pool
	prefix    string
}

var MAX_POOL_SIZE = 100

func newRedisPool(configSavePath string) (*redis.Pool, error) {
	var savePath string
	var poolsize int
	var password string
	var dbNum int
	var poollist *redis.Pool
	configs := strings.Split(configSavePath, ",")
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
	poollist = &redis.Pool{
		MaxIdle:     poolsize,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialTimeout(
				"tcp",
				savePath,
				time.Second,
				time.Second*12,
				time.Second,
			)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			_, err = c.Do("SELECT", dbNum)
			if err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	return poollist, poollist.Get().Err()
}
func NewRedisQueue(closeFunc *CloseFunc, config QueueStoreConfig) (QueueStoreInterface, error) {
	redisPool, err := newRedisPool(config.SavePath)
	if err != nil {
		return nil, err
	}
	result := &RedisQueueStore{
		redisPool: redisPool,
		prefix:    config.SavePrefix,
	}
	closeFunc.AddCloseHandler(func() {
		redisPool.Close()
	})
	return NewBasicQueue(result), nil
}

func (this *RedisQueueStore) Produce(topicId string, data interface{}) error {
	c := this.redisPool.Get()
	defer c.Close()

	_, err := c.Do("LPUSH", this.prefix+topicId, data)
	if err != nil {
		return err
	}
	return nil
}

func (this *RedisQueueStore) consumeData(topicId string, timeout int) (interface{}, error) {
	var topic interface{}
	var data interface{}

	c := this.redisPool.Get()
	defer c.Close()
	reply, err := redis.Values(c.Do("BRPOP", this.prefix+topicId, timeout))
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
	return data, nil
}

func (this *RedisQueueStore) Consume(topicId string, listener QueueListener) error {
	go func() {
		for {
			data, err := this.consumeData(topicId, 10)
			if err != nil {
				if strings.Index(err.Error(), "get on closed pool") != -1 {
					return
				} else {
					listener(err)
				}
			}
			if data == nil {
				continue
			} else {
				listener(data)
			}
		}
	}()
	return nil
}
