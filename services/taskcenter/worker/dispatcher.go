package worker

import (
	"aposervice/services/taskcenter/adapter"
	"fxlibraries/loggers"
	"time"
)

type Dispatcher struct {
	Interval time.Duration
}

func (dispatcher *Dispatcher) Start() {
	for {
		loggers.Info.Printf("Dispatcher start")
		dispatcher.dispatchTasks()
		loggers.Info.Printf("Dispatcher end")
		time.Sleep(dispatcher.Interval)
	}
}

func (dispatcher *Dispatcher) dispatchTasks() {
}
