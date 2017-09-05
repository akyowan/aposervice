package adapter

import (
	"aposervice/services/taskcenter/config"
	"fxlibraries/mongo"
	"fxlibraries/mysql"
)

var dbPool *mysql.DBPool
var mgoPool *mongo.MgoPool

func init() {
	dbPool = mysql.NewDBPool(mysql.DBPoolConfig{
		Host:         config.Conf.DBConf.Host,
		Port:         config.Conf.DBConf.Port,
		User:         config.Conf.DBConf.User,
		DBName:       config.Conf.DBConf.DBName,
		Password:     config.Conf.DBConf.Password,
		MaxIdleConns: 4,
		MaxOpenConns: 8,
		Debug:        config.IsDebug,
	})
	mgoPool = mongo.NewPool(&mongo.MongodbConfig{
		Host:   config.Conf.MongoConf.Host,
		Port:   config.Conf.MongoConf.Port,
		DBName: config.Conf.MongoConf.DBName,
	})

}
