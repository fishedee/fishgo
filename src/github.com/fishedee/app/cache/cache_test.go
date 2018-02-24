package web

import (
	"fmt"
	. "github.com/fishedee/assert"
	"testing"
	"time"
)

func newCacheForTest(t *testing.T, config CacheConfig) Cache {
	manager, err := NewCache(config)
	AssertEqual(t, err, nil, 0)
	return manager
}

func getExistData(t *testing.T, cache Cache, key string, index int) string {
	data, err := cache.Get(key)
	AssertEqual(t, err, nil, index)
	return data
}

func getNoExistData(t *testing.T, cache Cache, key string, index int) string {
	data, err := cache.Get(key)
	AssertEqual(t, err, nil, index)
	return data
}

func delData(t *testing.T, cache Cache, key string, index int) {
	err := cache.Del(key)
	AssertEqual(t, err, nil, index)
}

func setData(t *testing.T, cache Cache, key string, value string, duration time.Duration, index int) {
	err := cache.Set(key, value, duration)
	AssertEqual(t, err, nil, index)
}

func TestCacheBasic(t *testing.T) {
	testCaseDriver := []Cache{
		newCacheForTest(t, CacheConfig{
			Driver:     "memory",
			GcInterval: 1,
			SavePrefix: "cache:",
		}),
		newCacheForTest(t, CacheConfig{
			Driver:     "redis",
			SavePath:   "127.0.0.1:6379,100,13420693396",
			SavePrefix: "cache:",
		}),
	}
	for index, manager := range testCaseDriver {
		//清空数据
		delData(t, manager, "key1", index)
		delData(t, manager, "key2", index)
		delData(t, manager, "key3", index)
		delData(t, manager, "key4", index)

		//get与set
		setData(t, manager, "key1", "value1", time.Minute, index)
		setData(t, manager, "key2", "100", time.Minute, index)
		setData(t, manager, "key3", "value3", time.Minute, index)
		setData(t, manager, "key4", "", time.Minute, index)

		AssertEqual(t, getExistData(t, manager, "key1", index), "value1", index)
		AssertEqual(t, getExistData(t, manager, "key2", index), "100", index)
		AssertEqual(t, getExistData(t, manager, "key3", index), "value3", index)
		AssertEqual(t, getExistData(t, manager, "key4", index), "", index)

		//del
		delData(t, manager, "key3", index)
		AssertEqual(t, getExistData(t, manager, "key1", index), "value1", index)
		AssertEqual(t, getExistData(t, manager, "key2", index), "100", index)
		AssertEqual(t, getNoExistData(t, manager, "key3", index), "", index)

		//timeout expire
		setData(t, manager, "key2", "101", time.Second, index)
		time.Sleep(time.Second * 3)
		AssertEqual(t, getNoExistData(t, manager, "key2", index), "", index)
	}
}

func TestCacheMemoize(t *testing.T) {
	testCaseDriver := []Cache{
		newCacheForTest(t, CacheConfig{
			Driver:     "memory",
			GcInterval: 1,
			SavePrefix: "cache:",
		}),
		newCacheForTest(t, CacheConfig{
			Driver:     "redis",
			SavePath:   "127.0.0.1:6379,100,13420693396",
			SavePrefix: "cache:",
		}),
	}
	for index, manager := range testCaseDriver {
		//清空数据
		for i := 0; i != 100; i++ {
			delData(t, manager, fmt.Sprintf("%v", i), index)
		}
		//重置function
		var origin func(n int) int
		origin = func(n int) int {
			return manager.MustMemoize(fmt.Sprintf("%v", n), func() int {
				if n <= 2 {
					return 1
				}
				return origin(n-1) + origin(n-2)
			}, time.Minute).(int)
		}
		AssertEqual(t, origin(1), 1)
		AssertEqual(t, origin(2), 1)
		AssertEqual(t, origin(40), 102334155)
	}
}
