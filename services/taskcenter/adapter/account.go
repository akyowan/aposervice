package adapter

import "fxlibraries/loggers"

func GetCacheFree() (int64, error) {
	c := mgoPool.C("account_cache")
	maxCache := MaxApoCache()
	curCount, err := c.Count()
	if err != nil {
		loggers.Error.Printf("GetCacheFree get accounts cache count error:%s", err.Error())
		return 0, err
	}
	free := maxCache - int64(curCount)
	if free >= 0 {
		return free, nil
	}
	return 0, nil
}

func GetAppCacheCount(appID string) (int64, error) {
	//c := mgoPool.C("account_cache")
	return 0, nil
}

func DispatchAccountCache(appID string, count int) (int, error) {
	return 0, nil
}
