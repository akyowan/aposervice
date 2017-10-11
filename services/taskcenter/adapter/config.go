package adapter

import (
	"fxlibraries/loggers"
)

var (
	ConfigKey           string
	MaxApoCacheKey      string
	MaxApoCacheDefault  uint64
	MaxTaskCacheKey     string
	MaxTaskCacheDefault uint64
	MinTaskCacheKey     string
	MinTaskCacheDefault uint64
)

func init() {
	ConfigKey = "APOSERVICE.CONFIG"
	MaxApoCacheKey = "TASKCENTER.DISPATCHER.MAX_APO_CACHE"
	MaxApoCacheDefault = 20000
	MaxTaskCacheKey = "TASKCENTER.DISPACHTER.MAX_TASK_CACHE"
	MaxTaskCacheDefault = 1000
	MaxTaskCacheKey = "TASKCENTER.DISPACHTER.MIN_TASK_CACHE"
	MaxTaskCacheDefault = 100
}

func MaxTaskCache() uint64 {
	result := redisPool.HGet(ConfigKey, MaxTaskCacheKey)
	count, err := result.Uint64()
	if err != nil || count == 0 {
		loggers.Warn.Printf("MaxTaskCache get task max cache account from redis error:%s", err.Error())
		return MaxTaskCacheDefault
	}
	return count
}

func MinTaskCache() uint64 {
	result := redisPool.HGet(ConfigKey, MinTaskCacheKey)
	count, err := result.Uint64()
	if err != nil || count == 0 {
		loggers.Warn.Printf("MinTaskCache get task min cache account from redis error:%s", err.Error())
		return MinTaskCacheDefault
	}
	return count
}

func MaxApoCache() uint64 {
	result := redisPool.HGet(ConfigKey, MaxApoCacheKey)
	count, err := result.Uint64()
	if err != nil || count == 0 {
		loggers.Warn.Printf("MaxTaskCache get task max cache account from redis error:%s", err.Error())
		return MaxApoCacheDefault
	}
	return count
}
