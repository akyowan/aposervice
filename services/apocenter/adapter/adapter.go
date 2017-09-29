package adapter

import (
	"aposervice/config"
	"fxlibraries/mongo"
	"fxlibraries/redis"
	//"fxlibraries/mysql"
	"strconv"
)

//var dbPool *mysql.DBPool
var mgoPool *mongo.MgoPool
var redisPool *redis.RedisPool

func init() {
	//dbPool = mysql.NewPool(mysql.DBPoolConfig{
	//	Host:         config.Conf.Mysql.Host,
	//	Port:         config.Conf.Mysql.Port,
	//	User:         config.Conf.Mysql.User,
	//	DBName:       config.Conf.Mysql.DBName,
	//	Password:     config.Conf.Mysql.Password,
	//	MaxIdleConns: 4,
	//	MaxOpenConns: 8,
	//	Debug:        config.IsDebug,
	//})

	mgoPool = mongo.NewPool(&mongo.MongodbConfig{
		Host:   config.Conf.MongoDB.Host,
		Port:   config.Conf.MongoDB.Port,
		DBName: config.Conf.MongoDB.DBName,
		Debug:  config.IsDebug,
	})
	db, err := strconv.Atoi(config.Conf.Redis.DBName)
	if err != nil {
		panic(err)
	}
	redisPool = redis.NewPool(&redis.RedisConfig{
		Host: config.Conf.Redis.Host,
		DB:   db,
		Port: config.Conf.Redis.Port,
	})
}
