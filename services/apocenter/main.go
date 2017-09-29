package main

import (
	"aposervice/config"
	"aposervice/services/apocenter/app"
)

func main() {
	app.Start(config.Conf.Server.InternalListenAddress)
}
