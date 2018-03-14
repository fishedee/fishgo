package util

import (
	. "github.com/fishedee/assert"
	"testing"
	"time"
)

func newPoolTest(config *PoolConfig) *Pool {
	pool, err := NewPool(config)
	if err != nil {
		panic(err)
	}
	return pool
}

func TestPoolReuse(t *testing.T) {
	id := 1
	get := func() (interface{}, error) {
		oldId := id
		id++
		return oldId, nil
	}
	close := func(data interface{}) {

	}
	pools := []*Pool{
		newPoolTest(&PoolConfig{
			InitCap: 0,
			MaxCap:  10,
			Wait:    false,
			Get:     get,
			Close:   close,
		}),
		newPoolTest(&PoolConfig{
			InitCap: 0,
			MaxCap:  10,
			Wait:    true,
			Get:     get,
			Close:   close,
		}),
	}

	for _, pool := range pools {
		id = 1

		for i := 1; i <= 10; i++ {
			data, err := pool.Get()
			AssertEqual(t, err, nil)
			AssertEqual(t, data, i)
		}

		for i := 10; i >= 5; i-- {
			err := pool.Put(i, false)
			AssertEqual(t, err, nil)
		}

		for i := 10; i >= 5; i-- {
			data, err := pool.Get()
			AssertEqual(t, err, nil)
			AssertEqual(t, data, i)
		}

	}
}

func TestPoolNoWaitGet(t *testing.T) {
	id := 1

	get := func() (interface{}, error) {
		oldId := id
		id++
		return oldId, nil
	}
	close := func(data interface{}) {

	}

	pool := newPoolTest(&PoolConfig{
		InitCap: 0,
		MaxCap:  2,
		Wait:    false,
		Get:     get,
		Close:   close,
	})

	for i := 1; i <= 5; i++ {
		data, err := pool.Get()
		AssertEqual(t, err, nil)
		AssertEqual(t, data, i)
	}

	for i := 1; i <= 3; i++ {
		err := pool.Put(i, false)
		AssertEqual(t, err, nil)
	}

	for i := 1; i <= 2; i++ {
		data, err := pool.Get()
		AssertEqual(t, err, nil)
		AssertEqual(t, data, i)
	}

	for i := 1; i <= 2; i++ {
		data, err := pool.Get()
		AssertEqual(t, err, nil)
		AssertEqual(t, data, i+5)
	}
}

func TestPoolWaitGet(t *testing.T) {
	id := 1

	get := func() (interface{}, error) {
		oldId := id
		id++
		return oldId, nil
	}
	close := func(data interface{}) {

	}

	pool := newPoolTest(&PoolConfig{
		InitCap:     0,
		MaxCap:      2,
		Wait:        true,
		WaitTimeout: time.Second * 2,
		Get:         get,
		Close:       close,
	})

	for i := 1; i <= 2; i++ {
		data, err := pool.Get()
		AssertEqual(t, err, nil)
		AssertEqual(t, data, i)
	}

	go func() {
		time.Sleep(time.Second)
		err := pool.Put(1, false)
		AssertEqual(t, err, nil)
	}()

	//等待获取
	begin := time.Now().Unix()
	data, err := pool.Get()
	end := time.Now().Unix()
	AssertEqual(t, err, nil)
	AssertEqual(t, data, 1)
	AssertEqual(t, end-begin >= 1, true)

	//超时获取
	begin2 := time.Now().Unix()
	data2, err := pool.Get()
	end2 := time.Now().Unix()
	AssertEqual(t, err != nil, true)
	AssertEqual(t, data2, nil)
	AssertEqual(t, end2-begin2 >= 2, true)
}
