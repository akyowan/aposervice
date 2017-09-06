package app

import (
	"aposervice/services/taskcenter/worker"
	"fxlibraries/httpserver"
	"fxlibraries/loggers"
	"time"
)

func init() {
	loggers.Info.Printf("Initialize...\n")
}

func Start(addr string) {
	r := httpserver.NewRouter()
	loggers.Info.Printf("Starting TaskCenter External Service [\033[0;32;1mOK\t%+v\033[0m] \n", addr)

	// Start apo merger
	merger := worker.AppMerger{time.Second * 10}
	go merger.Start()

	panic(r.ListenAndServe(addr))
}
