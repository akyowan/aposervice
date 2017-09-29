package main

import (
	"aposervice/config"
	"aposervice/services/datacenter/app"
)

func main() {
	app.Start(config.Conf.Server.InternalListenAddress)
}
