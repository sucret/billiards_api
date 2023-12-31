// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameTaskLog = "task_log"

// TaskLog mapped from table <task_log>
type TaskLog struct {
	TaskLogID int64  `gorm:"column:task_log_id;type:bigint(20);primaryKey;autoIncrement:true" json:"task_log_id"`
	TaskID    int32  `gorm:"column:task_id;type:int(11);not null" json:"task_id"`
	Status    int    `gorm:"column:status;type:tinyint(4);not null" json:"status"` // 任务状态：1|执行中，2|执行成功，3|执行失败，4|手动取消
	StartTime Time   `gorm:"column:start_time;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"start_time"`
	EndTime   Time   `gorm:"column:end_time;type:timestamp" json:"end_time"`
	Log       string `gorm:"column:log;type:longtext" json:"log"`
	CreatedAt Time   `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt Time   `gorm:"column:updated_at;type:timestamp" json:"updated_at"`
}

// TableName TaskLog's table name
func (*TaskLog) TableName() string {
	return TableNameTaskLog
}
