package adapter

import (
	"fxlibraries/loggers"
)

var (
	ConfigKey           string
	MaxApoCacheKey      string
	MaxApoCacheDefault  int64
	MaxTaskCacheKey     string
	MaxTaskCacheDefault int64
)

func init() {
	ConfigKey = "APOSERVICE.CONFIG"
	MaxApoCacheKey = "TASKCENTER.DISPATCHER.MAX_APO_CACHE"
	MaxApoCacheDefault = 20000
	MaxTaskCacheKey = "TASKCENTER.DISPACHTER.MAX_TASK_CACHE"
	MaxTaskCacheDefault = 10000
}

func MaxTaskCache() int64 {
	result := redisPool.HGet(ConfigKey, MaxTaskCacheKey)
	count, err := result.Int64()
	if err != nil || count == 0 {
		loggers.Warn.Printf("MaxTaskCache get task max cache account from redis error:%s", err.Error())
		return MaxTaskCacheDefault
	}
	return count
}

func MaxApoCache() int64 {
	result := redisPool.HGet(ConfigKey, MaxApoCacheKey)
	count, err := result.Int64()
	if err != nil || count == 0 {
		loggers.Warn.Printf("MaxTaskCache get task max cache account from redis error:%s", err.Error())
		return MaxApoCacheDefault
	}
	return count
}
