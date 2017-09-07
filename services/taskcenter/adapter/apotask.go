package adapter

import (
	"aposervice/domain"
	"github.com/jinzhu/gorm"

	"gopkg.in/mgo.v2"

	"time"
)

func GetOnlineApoTasksFromDB() (map[int]domain.ApoTask, error) {
	db := dbPool.NewConn()
	var tasks []domain.ApoTask
	now := time.Now()
	db = db.Where("status = ?", domain.ApoTaskStatusStart)
	db = db.Where("start < ?", now).Where("end > ?", now)
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

func UpdateApoTasksToDB(tasks []domain.ApoTask) error {
	db := dbPool.NewConn().Begin()
	for i := range tasks {
		t := tasks[i]
		updates := map[string]interface{}{
			"doing_count":   t.DoingCount,
			"done_count":    t.DoneCount,
			"fail_count":    t.FailCount,
			"timeout_count": t.TimeoutCount,
		}
		if err := db.Model(&t).Updates(updates).Error; err != nil {
			db.Rollback()
			return err
		}
	}
	db.Commit()

	return nil
}

func MergeSubTaskToMain() error {
	db := dbPool.NewConn().Begin()
	now := time.Now()
	var tasks []domain.ApoSubTask
	dbResult := db.Where("status = ?", domain.ApoSubTaskEnable).Where("exec_time < ?", now).Find(&tasks)
	if dbResult.RecordNotFound() {
		db.Rollback()
		return nil
	}
	if dbResult.Error != nil {
		db.Rollback()
		return dbResult.Error
	}
	for i := range tasks {
		t := tasks[i]
		if err := db.Where("id = ?", t.ApoID).Update("total", gorm.Expr("total + ?", t.Count)).Error; err != nil {
			db.Rollback()
			return err
		}
		if err := db.Where("id = ?", t.ID).Update("status = ?", domain.ApoSubTaskDisable).Error; err != nil {
			db.Rollback()
			return err
		}
	}
	db.Commit()
	return nil
}

func UpdateApoTasksStatus() error {
	db := dbPool.NewConn()

	sCur := db.Where("status = ?", domain.ApoTaskStatusEnd).Where("total > (doing_count + done_count)")
	if err := sCur.Model(&domain.ApoTask{}).Update("status", domain.ApoTaskStatusStart).Error; err != nil {
		return err
	}

	now := time.Now()
	eCur := db.Where("status = ?", domain.ApoTaskStatusStart).Where("total <= done_count OR end <= ?", now)
	if err := eCur.Model(&domain.ApoTask{}).Update("Status", domain.ApoTaskStatusEnd).Error; err != nil {
		return err
	}

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
