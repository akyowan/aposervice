package adapter

import (
	"aposervice/domain"
	"github.com/jinzhu/gorm"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

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

func UpdateApoTasksToDB(tasks map[int]domain.ApoTask) error {
	db := dbPool.NewConn().Begin()
	for _, v := range tasks {
		updates := map[string]interface{}{
			"doing_count":   v.DoingCount,
			"done_count":    v.DoneCount,
			"fail_count":    v.FailCount,
			"timeout_count": v.TimeoutCount,
		}
		if err := db.Model(&v).Updates(updates).Error; err != nil {
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
	var mainTask domain.ApoTask
	for i := range tasks {
		t := tasks[i]
		if err := db.Table(mainTask.TableName()).Where("id = ?", t.ApoID).Update("total", gorm.Expr("total + ?", t.Count)).Error; err != nil {
			db.Rollback()
			return err
		}
		if err := db.Table(t.TableName()).Where("id = ?", t.ID).Update("status = ?", domain.ApoSubTaskDisable).Error; err != nil {
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

func GetApoTasksFromMongo() ([]domain.ApoTask, error) {
	c := mgoPool.C("apo_tasks")
	var tasks []domain.ApoTask
	err := c.Find(nil).Sort("start_time", "pub_time", "level").All(&tasks)
	if err == mgo.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	pre := tasks[0]
	for i := 1; i < len(tasks); i++ {
		cur := tasks[i]
		if cur.AppID != pre.AppID {
			pre = cur
			continue
		}
		for j := i + 1; j < len(tasks); j++ {
			apo := tasks[j]
			if apo.AppID != pre.AppID {
				tasks[i] = tasks[j]
				pre = cur
				break
			}
		}
		break
	}
	return tasks, nil
}

func GetApoTaskFromMongo(id int) (*domain.ApoTask, error) {
	c := mgoPool.C("apo_tasks")
	var task domain.ApoTask
	if err := c.FindId(id).One(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func GetFirstApoTask() (*domain.ApoTask, error) {
	c := mgoPool.C("apo_tasks")
	var task domain.ApoTask
	if err := c.Find(nil).Sort("update_time", "start_time", "end_time", "level").Limit(1).One(&task); err != nil {
		return nil, err
	}
	return nil, nil
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

func SaveTaskToMongo(task *domain.ApoTask) error {
	c := mgoPool.C("apo_tasks")
	data := domain.ApoTask{
		AppID:        task.AppID,
		AppName:      task.AppName,
		BundleID:     task.BundleID,
		Level:        task.Level,
		Total:        task.Total,
		RealTotal:    task.RealTotal,
		ApoKey:       task.ApoKey,
		AccountBrief: task.AccountBrief,
		Cycle:        task.Cycle,
		RemindCycle:  task.RemindCycle,
		UncatchDay:   task.UncatchDay,
		TypeModelID:  task.TypeModelID,
		AmoutModelID: task.AmoutModelID,
		PreaddCount:  task.PreaddCount,
		PreaddTime:   task.PreaddTime,
		StartTime:    task.StartTime,
		EndTime:      task.EndTime,
		CreatedAt:    task.CreatedAt,
		UpdatedAt:    task.UpdatedAt,
	}
	var err error
	if err = c.UpdateId(task.ID, bson.M{"$set": data}); err == nil {
		return nil
	}
	if err = c.Insert(&task); err == nil {
		return nil
	}

	return err
}

func DeleteApoTaskFromMongo(id int) error {
	c := mgoPool.C("apo_tasks")
	if err := c.RemoveId(id); err != nil {
		return err
	}
	return nil
}
