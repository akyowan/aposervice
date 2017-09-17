package domain

import "time"

type ApoTask struct {
	ID           int           `bson:"_id,omitempty" gorm:"primary_key;column:id;unique_index:devices_pkey"`
	AppID        string        `bson:"app_id,omitempty"`
	AppName      string        `bson:"app_name,omitempty"`
	BundleID     string        `bson:"bundle_id,omitempty"`
	Level        int           `bson:"level,omitempty"`
	Total        int           `bson:"total,omitempty"`
	RealTotal    int           `bson:"real_total,omitempty"`
	DoneCount    int           `bson:"done_count,omitempty"`
	DoingCount   int           `bson:"doing_count,omitempty"`
	TimeoutCount int           `bson:"timeount_count,omitempty"`
	FailCount    int           `bson:"fail_count,omitempty"`
	ApoKey       string        `bson:"apo_key,omitempty"`
	AccountBrief string        `bson:"account_brief,omitempty"`
	Cycle        int           `bson:"cycel,omitempty"`
	RemindCycle  int           `bson:"remind_cycle,omitempty"`
	UncatchDay   int           `bson:"uncatch_day,omitempty"`
	TypeModelID  int64         `bson:"type_model_id,omitempty"`
	AmoutModelID int64         `bson:"amount_model_id,omitempty"`
	Status       ApoTaskStatus `bson:"status,omitempty"`
	PreaddCount  int           `bson:"preadd_count,omitempty"`
	PreaddTime   *time.Time    `bson:"preadd_time,omitempty"`
	StartTime    *time.Time    `bson:"start_time,omitempty" gorm:"column:start"`
	EndTime      *time.Time    `bson:"end_time,omitempty" gorm:"column:end"`
	CreatedAt    *time.Time    `bson:"create_time,omitempty" gorm:"column:create_time"`
	UpdatedAt    *time.Time    `bson:"update_time,omitempty" gorm:"column:update_time"`
}

func (*ApoTask) TableName() string {
	return "apo_tasks"
}

type ApoTaskStatus int

const (
	_                    ApoTaskStatus = iota
	ApoTaskStatusStart                 // 开始
	ApoTaskStatusPause                 // 暂停
	ApoTaskStatusEnd                   // 完成
	ApoTaskStatusDeleted               // 删除
)

type ApoSubTask struct {
	ID       int        `json:"-" gorm:"primary_key;column:id;unique_index:devices_pkey"`
	ApoID    int        `json:"apo_id"`
	Count    int        `json:"count"`
	Status   int        `json:"status"`
	ExecTime *time.Time `json:"exec_time"`
	CreateAt *time.Time `json:"create_time" gorm:"column:create_time"`
}

func (*ApoSubTask) TableName() string {
	return "apo_sub_task"
}

type ApoSubTaskStatus int

const (
	ApoSubTaskDisable ApoSubTaskStatus = iota //  禁用
	ApoSubTaskEnable                          // 可用
)
