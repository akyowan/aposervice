package worker

import (
	"aposervice/domain"
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

	var mergedTasks []domain.ApoTask
	var deletedTasks []domain.ApoTask

	if dbTasks != nil {
		for _, v := range dbTasks {
			id := v.ID
			if t, ok := cacheTasks[id]; ok {
				task := merger.mergeTask(&v, &t)
				mergedTasks = append(mergedTasks, *task)
			} else {
				mergedTasks = append(mergedTasks, v)
			}
		}
	}

	if cacheTasks != nil {
		for _, v := range cacheTasks {
			id := v.ID
			if _, ok := dbTasks[id]; !ok {
				if err := adapter.DeleteApoTaskFromMongo(id); err != nil {
					loggers.Error.Printf("AppMerger delete task:%d from mongo error:%s", id, err.Error())
					continue
				}
				deletedTasks = append(deletedTasks, v)
			}
		}
	}

	if err := adapter.UpdateApoTasksToDB(mergedTasks); err != nil {
		loggers.Error.Printf("AppMerger update merged apo tasks to db error:%s", err.Error())
	}
	if err := adapter.SaveTasksToMongo(mergedTasks); err != nil {
		loggers.Error.Printf("AppMerger save apo tasks to mongo error:%s", err.Error())
	}

	if err := adapter.UpdateApoTasksToDB(deletedTasks); err != nil {
		loggers.Error.Printf("AppMerger update deleted apo tasks to db error:%s", err.Error())
	}
}

func (merger *AppMerger) mergeTask(dbTask *domain.ApoTask, cacheTask *domain.ApoTask) *domain.ApoTask {
	(*dbTask).DoingCount = (*cacheTask).DoingCount
	(*dbTask).DoneCount = (*cacheTask).DoneCount
	(*dbTask).TimeoutCount = (*cacheTask).TimeoutCount
	(*dbTask).FailCount = (*cacheTask).FailCount
	return dbTask
}
