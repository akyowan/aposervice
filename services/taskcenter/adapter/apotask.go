package adapter

import (
	"aposervice/domain"
	"time"
)

func GetOnlineApoTasksFromDB() (map[int]domain.ApoTask, error) {
	db := dbPool.NewConn()
	var tasks []domain.ApoTask
	now := time.Now()
	db = db.Where("status = ?", domain.ApoTaskStatusStart)
	db = db.Where("start_time < ?", now).Where("end_time < ?", now)
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
		if err := db.Model(&t).Updates(&updates).Error; err != nil {
			db.Rollback()
			return err
		}
	}

	db.Commit()
	return nil
}

func GetAllApoTaskFromMongo() (map[int]domain.ApoTask, error) {
	return nil, nil
}

func SaveTasksToMongo(tasks []domain.ApoTask) error {
	return nil
}

func GetFirstApoTask() (*domain.ApoTask, error) {
	return nil, nil
}

func DeleteApoTaskFromMongo(id int) error {
	return nil
}
