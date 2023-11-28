// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameTask = "task"

// Task mapped from table <task>
type Task struct {
	TaskID              int32  `gorm:"column:task_id;type:int(10) unsigned;primaryKey;autoIncrement:true" json:"task_id"`     // 主键
	Name                string `gorm:"column:name;type:varchar(64);not null" json:"name"`                                     // 任务名称
	Spec                string `gorm:"column:spec;type:varchar(64);not null" json:"spec"`                                     // crontab 表达式
	Command             string `gorm:"column:command;type:varchar(255);not null" json:"command"`                              // 执行命令
	ProcessNum          int32  `gorm:"column:process_num;type:int(11);not null;default:1" json:"process_num"`                 // 进程数（单机）
	Timeout             int32  `gorm:"column:timeout;type:int(10) unsigned;not null;default:60" json:"timeout"`               // 超时时间(单位:秒)
	RetryTimes          int32  `gorm:"column:retry_times;type:tinyint(4);not null;default:3" json:"retry_times"`              // 重试次数
	RetryInterval       int32  `gorm:"column:retry_interval;type:int(11);not null;default:60" json:"retry_interval"`          // 重试间隔(单位:秒)
	NotifyStatus        int32  `gorm:"column:notify_status;type:tinyint(4);not null" json:"notify_status"`                    // 执行结束是否通知 1:不通知 2:失败通知 3:结束通知 4:结果关键字匹配通知
	NotifyType          int32  `gorm:"column:notify_type;type:tinyint(4);not null;default:1" json:"notify_type"`              // 通知类型 1:邮件 2:webhook
	NotifyReceiverEmail string `gorm:"column:notify_receiver_email;type:varchar(255);not null" json:"notify_receiver_email"`  // 通知者邮箱地址(多个用,分割)
	NotifyKeyword       string `gorm:"column:notify_keyword;type:varchar(255);not null" json:"notify_keyword"`                // 通知匹配关键字(多个用,分割)
	Remark              string `gorm:"column:remark;type:varchar(100);not null" json:"remark"`                                // 备注
	Status              int32  `gorm:"column:status;type:tinyint(4);not null;default:1" json:"status"`                        // 状态1:启用  2:停用
	CreatedAt           Time   `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	UpdatedAt           Time   `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
}

// TableName Task's table name
func (*Task) TableName() string {
	return TableNameTask
}
