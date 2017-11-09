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

	// Delete comments
	r.RouteHandleFunc("/comments", handler.DeleteComments).Queries("action", "delete").Methods("POST")
	r.RouteHandleFunc("/comments/{appID}", handler.DeleteAppComments).Queries("type", "app_id").Methods("DELETE")
	r.RouteHandleFunc("/comments/{id}", handler.DeleteComment).Methods("DELETE")

	// Post comments
	r.RouteHandleFunc("/comments", handler.AddComments).Methods("POST")

	// Get comments
	r.RouteHandleFunc("/comments", handler.GetComments).Methods("GET")

	// Update comment
	r.RouteHandleFunc("/comments/{id}", handler.UpdateComment).Methods("PATCH")
	r.RouteHandleFunc("/comments", handler.UpdateAndDeleteComments).Methods("PATCH")

	loggers.Info.Printf("Starting data center external service [\033[0;32;1mOK\t%+v\033[0m] \n", addr)
	panic(r.ListenAndServe(addr))
}
