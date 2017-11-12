package worker

import (
	"aposervice/services/apocenter/adapter"
	"fxlibraries/loggers"
	"time"
)

func RecycleCommentRun() {
	for {
		if err := adapter.RecycleComment(); err != nil {
			loggers.Error.Printf("RecycleCommentRun error:%s", err.Error())
		}
		time.Sleep(time.Second * 30)
	}
}
