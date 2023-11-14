package service

import (
	"gin-api/pkg/mysql"
	"gin-api/pkg/mysql/model"
	"time"

	"gorm.io/gorm"
)

type taskLog struct {
	db *gorm.DB
}

var TaskLogService = &taskLog{
	db: mysql.GetDB(),
}

// 任务开始执行时创建日志
func (l *taskLog) SaveLog(taskID int32) (taskLog *model.TaskLog) {
	taskLog = &model.TaskLog{
		TaskID:    taskID,
		Status:    1,
		StartTime: model.Time(time.Now()),
		CreatedAt: model.Time(time.Now()),
	}

	l.db.Create(&taskLog)
	return
}

func (l *taskLog) UpdateLog(taskLog *model.TaskLog) {
	l.db.Save(&taskLog)
}
