package main

import (
	"aposervice/config"
	"aposervice/services/taskcenter/app"
)

func main() {
	app.Start(config.Conf.Server.InternalListenAddress)
}
