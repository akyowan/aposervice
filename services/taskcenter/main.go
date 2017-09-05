package main

import (
	"aposervice/services/taskcenter/app"
	"aposervice/services/taskcenter/config"
)

func main() {
	app.Start(config.Conf.ServerConf.InternalListenAddress)
}
