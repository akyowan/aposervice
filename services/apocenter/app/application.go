package app

import (
	"aposervice/services/apocenter/handler"
	"aposervice/services/apocenter/worker"
	"fxlibraries/httpserver"
	"fxlibraries/loggers"
)

func init() {
	loggers.Info.Printf("Initialize...\n")
}

func Start(addr string) {
	r := httpserver.NewRouter()
	loggers.Info.Printf("Starting ApoCenter External Service [\033[0;32;1mOK\t%+v\033[0m] \n", addr)
	r.RouteHandleFunc("/comments/{appID}", handler.GetAppComments).Methods("GET")
	r.RouteHandleFunc("/comments", handler.ReportComment).Queries("action", "report").Methods("POST")
	go worker.RecycleCommentRun()
	panic(r.ListenAndServe(addr))
}
