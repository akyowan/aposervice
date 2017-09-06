package adapter

import (
	"aposervice/domain"
	"fxlibraries/loggers"

	"gopkg.in/mgo.v2"

	"time"
)

func GetOnlineApoTasksFromDB() (map[int]domain.ApoTask, error) {
	db := dbPool.NewConn()
	var tasks []domain.ApoTask
	now := time.Now()
	db = db.Where("status = ?", domain.ApoTaskStatusStart)
	db = db.Where("start_time < ?", now).Where("end_time > ?", now)
	dbResult := db.Find(&tasks)
	if dbResult.RecordNotFound() {
		return nil, nil
	}
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}

	tasksMap := make(map[int]domain.ApoTask)
	for i := range tasks {
		tasksMap[tasks[i].ID] = tasks[i]
	}
	return tasksMap, nil
}

func UpdateApoTasksToDB(tasks []domain.ApoTask, withStatus bool) error {
	db := dbPool.NewConn().Begin()
	for i := range tasks {
		t := tasks[i]
		updates := map[string]interface{}{
			"doing_count":   t.DoingCount,
			"done_count":    t.DoneCount,
			"fail_count":    t.FailCount,
			"timeout_count": t.TimeoutCount,
		}
		if withStatus {
			updates["status"] = t.Status
		}
		loggers.Info.Printf("update %d to db", t.ID)
		if err := db.Model(&t).Updates(updates).Error; err != nil {
			db.Rollback()
			return err
		}
	}
	db.Commit()

	return nil
}

func GetAllApoTaskFromMongo() (map[int]domain.ApoTask, error) {
	c := mgoPool.C("apo_tasks")
	var tasks []domain.ApoTask
	err := c.Find(nil).All(&tasks)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	tasksMap := make(map[int]domain.ApoTask)
	for i := range tasks {
		tasksMap[tasks[i].ID] = tasks[i]
	}
	return tasksMap, nil
}

func GetApoTaskFromMongo(id int) (*domain.ApoTask, error) {
	c := mgoPool.C("apo_tasks")
	var task domain.ApoTask
	if err := c.FindId(id).One(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func SaveTasksToMongo(tasks []domain.ApoTask) error {
	c := mgoPool.C("apo_tasks")
	for i := range tasks {
		t := tasks[i]
		if _, err := c.UpsertId(t.ID, &t); err != nil {
			return err
		}
	}

	return nil
}

func GetFirstApoTask() (*domain.ApoTask, error) {
	c := mgoPool.C("apo_tasks")
	var task domain.ApoTask
	if err := c.Find(nil).Sort("update_time", "start_time", "end_time", "level").Limit(1).One(&task); err != nil {
		return nil, err
	}
	return nil, nil
}

func DeleteApoTaskFromMongo(id int) error {
	c := mgoPool.C("apo_tasks")
	if err := c.RemoveId(id); err != nil {
		return err
	}
	return nil
}
