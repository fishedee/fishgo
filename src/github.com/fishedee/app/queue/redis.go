package queue

import (
	"encoding/json"
	. "github.com/fishedee/app/log"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/redsync.v1"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type redisConfig struct {
	savePath string
	password string
	dbNum    int
	poolSize int
}

type redisQueueChannel struct {
	conn redis.Conn
}

type redisQueueStore struct {
	redisPool        *redis.Pool
	log              Log
	redisConfig      redisConfig
	config           QueueConfig
	waitgroup        *sync.WaitGroup
	isClose          int32
	consumeListeners sync.Map
	router           *map[string][]string
	exitChan         chan bool
	closeChan        chan bool
}

var MAX_POOL_SIZE = 100

func parseRedisConfig(savepath string) redisConfig {
	var savePath string
	var password string
	var dbNum int
	var poolsize int
	configs := strings.Split(savepath, ",")
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
	return redisConfig{
		savePath: savePath,
		password: password,
		dbNum:    dbNum,
		poolSize: poolsize,
	}
}
func newRedisQueue(log Log, config QueueConfig) (queueStoreInterface, error) {
	redisConfig := parseRedisConfig(config.SavePath)
	if config.RetryInterval == 0 {
		config.RetryInterval = 5
	}

	result := &redisQueueStore{
		config:      config,
		redisConfig: redisConfig,
		log:         log,
		waitgroup:   &sync.WaitGroup{},
		router:      &map[string][]string{},
		closeChan:   make(chan bool, 16),
		exitChan:    make(chan bool, 16),
	}
	var err error
	result.redisPool, err = result.getConnectPool()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (this *redisQueueStore) getConnectPool() (*redis.Pool, error) {
	poollist := &redis.Pool{
		MaxIdle:     this.redisConfig.poolSize,
		IdleTimeout: 240 * time.Second,
		Dial:        this.getConnect,
	}
	return poollist, poollist.Get().Err()
}

func (this *redisQueueStore) getConnect() (redis.Conn, error) {
	c, err := redis.DialTimeout(
		"tcp",
		this.redisConfig.savePath,
		time.Second,
		time.Second*12,
		time.Second,
	)
	if err != nil {
		return nil, err
	}
	if this.redisConfig.password != "" {
		if _, err := c.Do("AUTH", this.redisConfig.password); err != nil {
			c.Close()
			return nil, err
		}
	}
	_, err = c.Do("SELECT", this.redisConfig.dbNum)
	if err != nil {
		c.Close()
		return nil, err
	}
	return c, err
}

func (this *redisQueueStore) Produce(topicId string, data []byte) error {
	c := this.redisPool.Get()
	defer c.Close()

	router := this.getRouter()
	topicRouter, isExist := router[topicId]
	if isExist == false {
		return nil
	}

	var lastErr error
	for _, singleQueue := range topicRouter {
		_, err := c.Do("LPUSH", singleQueue, data)
		if err != nil {
			lastErr = err
		}
	}

	return lastErr
}

func (this *redisQueueStore) consumeData(connect redis.Conn, queueName string, timeout int) ([]byte, error) {
	var topic interface{}
	var data interface{}

	reply, err := redis.Values(connect.Do("BRPOP", this.config.SavePrefix+queueName, timeout))
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

func (this *redisQueueStore) singleConsume(queueName string, listener queueStoreListener) error {
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
		data, err := this.consumeData(conn, queueName, 5)
		isExit := atomic.LoadInt32(&this.isClose)
		if err != nil {
			if isExit == 1 {
				return nil
			} else {
				return err
			}
		}
		if data == nil {
			if isExit == 1 {
				return nil
			} else {
				continue
			}
		}
		listener(data)
	}
}

func (this *redisQueueStore) Consume(topicId string, queueName string, poolSize int, listener queueStoreListener) error {
	err := this.setRedisRouter(topicId, this.config.SavePrefix+queueName)
	if err != nil {
		return err
	}
	for i := 0; i < poolSize; i++ {
		this.waitgroup.Add(1)
		go func() {
			for {
				err := this.singleConsume(queueName, listener)
				if err != nil {
					this.log.Critical("Queue Redis consume error :%v, will be retry in %v seconds", err, this.config.RetryInterval)
					time.Sleep(time.Duration(int(time.Second) * this.config.RetryInterval))
				} else {
					this.waitgroup.Done()
					break
				}
			}
		}()
	}
	return nil
}

func (this *redisQueueStore) getRouter() map[string][]string {
	router := *(*map[string][]string)(atomic.LoadPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&this.router)),
	))
	return router
}

func (this *redisQueueStore) setRouter(data map[string][]string) {
	atomic.StorePointer(
		(*unsafe.Pointer)(unsafe.Pointer(&this.router)),
		unsafe.Pointer(&data),
	)
}

func (this *redisQueueStore) updateRouter() {
	result, err := this.getRedisRouter()
	if err != nil {
		this.log.Critical("getRedisRouter fail: %v", err)
		return
	}
	this.setRouter(result)
}

func (this *redisQueueStore) getRedisRouter() (map[string][]string, error) {
	c := this.redisPool.Get()
	defer c.Close()
	valueGet, err := redis.Bytes(c.Do("GET", this.config.SavePrefix+"queue_topic_info"))
	if err == redis.ErrNil {
		return map[string][]string{}, nil
	}
	if err != nil {
		return nil, err
	}
	result := map[string][]string{}
	err = json.Unmarshal(valueGet, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type redisCurrentPool struct {
	conn redis.Conn
}

func (this *redisCurrentPool) Get() redis.Conn {
	return this.conn
}

func (this *redisQueueStore) setRedisRouter(topicId string, queueName string) error {
	sync := redsync.New([]redsync.Pool{this.redisPool})
	mutex := sync.NewMutex(this.config.SavePrefix + "queue_topic_mutex")
	err := mutex.Lock()
	if err != nil {
		return err
	}
	defer mutex.Unlock()

	router, err := this.getRedisRouter()
	if err != nil {
		return err
	}
	topicInfo := router[topicId]
	hasExist := false
	for _, singleQueue := range topicInfo {
		if singleQueue == queueName {
			hasExist = true
			break
		}
	}
	if hasExist == false {
		router[topicId] = append(topicInfo, queueName)
	}

	data, err := json.Marshal(router)
	if err != nil {
		return err
	}
	c := this.redisPool.Get()
	defer c.Close()
	_, err = c.Do("SET", this.config.SavePrefix+"queue_topic_info", data)
	if err != nil {
		return err
	}
	return nil
}

func (this *redisQueueStore) Run() error {
	isRun := true

	this.updateRouter()
	tick := time.NewTicker(time.Second)
	for isRun {
		select {
		case <-this.closeChan:
			isRun = false
			break
		case <-tick.C:
			this.updateRouter()
			break
		}
	}

	this.waitgroup.Wait()
	this.exitChan <- true
	return nil
}

func (this *redisQueueStore) Close() {
	this.closeChan <- true

	atomic.StoreInt32(&this.isClose, 1)
	this.redisPool.Close()
	this.closeListener()

	<-this.exitChan
}

func (this *redisQueueStore) closeListener() {
	this.consumeListeners.Range(func(key, value interface{}) bool {
		conn := key.(redis.Conn)
		conn.Close()
		return true
	})
}
