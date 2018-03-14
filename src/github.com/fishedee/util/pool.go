package util

import (
	"errors"
	"sync"
	"time"
)

type Pool struct {
	mutex       sync.RWMutex
	connect     chan interface{}
	get         func() (interface{}, error)
	close       func(interface{})
	wait        bool
	waitTimeout time.Duration
}

type PoolConfig struct {
	InitCap     int
	MaxCap      int
	Wait        bool
	WaitTimeout time.Duration
	Get         func() (interface{}, error)
	Close       func(interface{})
}

func NewPool(config *PoolConfig) (*Pool, error) {
	if config.MaxCap <= 0 || config.InitCap < 0 {
		return nil, errors.New("unset initCap or maxCap")
	}
	if config.InitCap > config.MaxCap {
		return nil, errors.New("initcap must less than maxCap")
	}
	if config.WaitTimeout == 0 {
		config.WaitTimeout = time.Second
	}
	pool := &Pool{
		mutex:       sync.RWMutex{},
		connect:     make(chan interface{}, config.MaxCap),
		get:         config.Get,
		close:       config.Close,
		wait:        config.Wait,
		waitTimeout: config.WaitTimeout,
	}

	for i := 0; i < config.InitCap; i++ {
		conn, err := pool.get()
		if err != nil {
			pool.close(conn)
			return nil, err
		}
		pool.connect <- conn
	}
	if config.Wait {
		for i := 0; i < config.MaxCap-config.InitCap; i++ {
			pool.connect <- nil
		}
	}
	return pool, nil
}

func (this *Pool) Get() (interface{}, error) {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	conns := this.connect

	if conns == nil {
		return nil, errors.New("pool has closed")
	}
	if this.wait {
		select {
		case conn := <-conns:
			if conn == nil {
				var err error
				conn, err = this.get()
				if err != nil {
					conns <- nil
					return nil, err
				}
			}
			return conn, nil
		case <-time.After(this.waitTimeout):
			return nil, errors.New("pool get timeout")
		}
	} else {
		select {
		case conn := <-conns:
			return conn, nil
		default:
			conn, err := this.get()
			return conn, err
		}
	}

}

func (this *Pool) Put(conn interface{}, forceClose bool) error {
	if conn == nil {
		return errors.New("invalid conn put,it is nil")
	}
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	conns := this.connect

	if conns == nil {
		this.close(conn)
		return nil
	}

	if this.wait {
		if forceClose {
			this.close(conn)
			conn = nil
		}
		conns <- conn
		return nil
	} else {
		if forceClose {
			this.close(conn)
			return nil
		}
		select {
		case conns <- conn:
			return nil
		default:
			this.close(conn)
			return nil
		}
	}

}

func (this *Pool) Close() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	conns := this.connect
	if conns == nil {
		return
	}
	close(conns)
	for conn := range conns {
		this.close(conn)
	}
	this.connect = nil
}
