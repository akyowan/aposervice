package main

import (
	"aposervice/services/datacenter/app"
	"aposervice/services/datacenter/config"
)

func main() {
	app.Start(config.Conf.Server.InternalListenAddress)
}
