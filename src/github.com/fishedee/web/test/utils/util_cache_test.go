package web

import (
	. "github.com/fishedee/web"
	"reflect"
	"testing"
	"time"
)

func assertCacheEqual(t *testing.T, left interface{}, right interface{}, index int) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("case :%v ,%+v != %+v", index, left, right)
	}
}

func newCacheForTest(t *testing.T, config CacheConfig) Cache {
	manager, err := NewCache(config)
	assertCacheEqual(t, err, nil, 0)
	return manager
}

func getExistData(t *testing.T, cache Cache, key string, index int) string {
	data, isOk := cache.Get(key)
	assertCacheEqual(t, isOk, true, index)
	return data
}

func getNoExistData(t *testing.T, cache Cache, key string, index int) string {
	data, isOk := cache.Get(key)
	assertCacheEqual(t, isOk, false, index)
	return data
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
		manager.Del("key1")
		manager.Del("key2")
		manager.Del("key3")

		//get与set
		manager.Set("key1", "value1", time.Minute)
		manager.Set("key2", "100", time.Minute)
		manager.Set("key3", "value3", time.Minute)
		manager.Set("key4", "", time.Minute)
		assertCacheEqual(t, getExistData(t, manager, "key1", index), "value1", index)
		assertCacheEqual(t, getExistData(t, manager, "key2", index), "100", index)
		assertCacheEqual(t, getExistData(t, manager, "key3", index), "value3", index)
		assertCacheEqual(t, getExistData(t, manager, "key4", index), "", index)

		//del
		manager.Del("key3")
		assertCacheEqual(t, getExistData(t, manager, "key1", index), "value1", index)
		assertCacheEqual(t, getExistData(t, manager, "key2", index), "100", index)
		assertCacheEqual(t, getNoExistData(t, manager, "key3", index), "", index)

		//timeout expire
		manager.Set("key2", "101", time.Second)
		time.Sleep(time.Second * 3)
		assertCacheEqual(t, getNoExistData(t, manager, "key2", index), "", index)
	}
}
