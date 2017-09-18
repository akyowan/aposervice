package app

import (
	"aposervice/services/datacenter/handler"
	"fxlibraries/httpserver"
	"fxlibraries/loggers"
)

func init() {
	loggers.Info.Printf("Initialize...\n")
}

func Start(addr string) {
	r := httpserver.NewRouter()
	r.RouteHandleFunc("/comments", handler.AddComments).Methods("POST")
	r.RouteHandleFunc("/comments", handler.GetComments).Methods("GET")

	loggers.Info.Printf("Starting TaskCenter External Service [\033[0;32;1mOK\t%+v\033[0m] \n", addr)
	panic(r.ListenAndServe(addr))
}
