package util

import (
	"reflect"
	"testing"
	"time"
)

func assertCacheEqual(t *testing.T, left interface{}, right interface{}, index int) {
	if reflect.DeepEqual(left, right) == false {
		t.Errorf("case :%v ,%+v != %+v", index, left, right)
	}
}

func newCacheManagerForTest(t *testing.T, config CacheManagerConfig) *CacheManager {
	manager, err := newCacheManager(config)
	assertCacheEqual(t, err, nil, 0)
	return manager
}

func TestCacheBasic(t *testing.T) {
	testCaseDriver := []*CacheManager{
		newCacheManagerForTest(t, CacheManagerConfig{
			Driver:     "memory",
			GcInterval: 1,
			SavePrefix: "cache:",
		}),
		newCacheManagerForTest(t, CacheManagerConfig{
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
		assertCacheEqual(t, manager.Get("key1"), "value1", index)
		assertCacheEqual(t, manager.Get("key2"), "100", index)
		assertCacheEqual(t, manager.Get("key3"), "value3", index)

		//del
		manager.Del("key3")
		assertCacheEqual(t, manager.Get("key1"), "value1", index)
		assertCacheEqual(t, manager.Get("key2"), "100", index)
		assertCacheEqual(t, manager.Get("key3"), "", index)

		//timeout expire
		manager.Set("key2", "101", time.Second)
		time.Sleep(time.Second * 3)
		assertCacheEqual(t, manager.Get("key2"), "", index)
	}
}
