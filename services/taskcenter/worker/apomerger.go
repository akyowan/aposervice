package worker

import (
	"aposervice/services/taskcenter/adapter"
	"fxlibraries/loggers"
	"time"
)

type AppMerger struct {
	Interval time.Duration
}

func (merger *AppMerger) Start() {
	for {
		loggers.Info.Printf("AppMerger start")
		merger.mergeTasks()
		loggers.Info.Printf("AppMerger end")
		time.Sleep(merger.Interval)
	}
}

func (merger *AppMerger) mergeTasks() {
	if err := adapter.MergeSubTaskToMain(); err != nil {
		loggers.Error.Printf("AppMerger merge sub task to main task error:%s", err.Error())
		return
	}
	if err := adapter.UpdateApoTasksStatus(); err != nil {
		loggers.Error.Printf("AppMerger refresh apo tasks status error:%s", err.Error())
		return
	}
	dbTasks, err := adapter.GetOnlineApoTasksFromDB()
	if err != nil {
		loggers.Error.Printf("AppMerger get online tasks from db error:%s", err.Error())
		return
	}

	cacheTasks, err := adapter.GetAllApoTaskFromMongo()
	if err != nil {
		loggers.Error.Printf("AppMerger get online tasks from cache error:%s", err.Error())
		return
	}

	for _, v := range dbTasks {
		if err := adapter.SaveTaskToMongo(&v); err != nil {
			loggers.Error.Printf("AppMerger save task to mongo error:%s", err.Error())
			return
		}
	}

	for _, v := range cacheTasks {
		id := v.ID
		if _, ok := dbTasks[id]; !ok {
			if err := adapter.DeleteApoTaskFromMongo(id); err != nil {
				loggers.Error.Printf("AppMerger delete task:%d from mongo error:%s", id, err.Error())
				continue
			}
		}
	}

	if err := adapter.UpdateApoTasksToDB(cacheTasks); err != nil {
		loggers.Error.Printf("AppMerger update deleted apo tasks to db error:%s", err.Error())
	}
}
