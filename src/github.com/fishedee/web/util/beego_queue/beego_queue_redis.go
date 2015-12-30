package beego_queue

import (
	"github.com/garyburd/redigo/redis"
	"strconv"
	"strings"
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
		poolsize, err := strconv.Atoi(configs[1])
		if err != nil || poolsize <= 0 {
			poolsize = MAX_POOL_SIZE
		} else {
			poolsize = poolsize
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
	poollist = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", savePath)
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
	}, poolsize)

	return poollist, poollist.Get().Err()
}
func NewRedisQueue(config BeegoQueueStoreConfig) (*RedisQueueStore, error) {
	redisPool, err := newRedisPool(config.SavePath)
	if err != nil {
		return nil, err
	}
	result := &RedisQueueStore{
		redisPool: redisPool,
		prefix:    config.SavePrefix,
	}
	return result, nil
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

func (this *RedisQueueStore) Consume(topicId string, listener BeegoQueueListener) error {
	go func() {
		for {
			data, err := this.consumeData(topicId, 30)
			if err != nil {
				listener(err)
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

func (this *RedisQueueStore) Publish(topicId string, data interface{}) error {
	c := this.redisPool.Get()
	defer c.Close()

	_, err := c.Do("PUBLISH", this.prefix+topicId, data)
	if err != nil {
		return err
	}
	return nil
}

func (this *RedisQueueStore) Subscribe(topicId string, listener BeegoQueueListener) error {
	c := this.redisPool.Get()

	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe(this.prefix + topicId)
	go func() {
		defer c.Close()
		for {
			switch n := psc.Receive().(type) {
			case redis.Message:
				listener(n.Data)
			case error:
				listener(n)
				return
			default:
				break
			}
		}
	}()
	return nil
}
