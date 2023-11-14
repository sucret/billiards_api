package request

type ChangeTaskStatus struct {
	TaskId int32 `form:"task_id" json:"task_id" binding:"required"`
	Status int32 `form:"status" json:"status" binding:"required,oneof=1 2"`
}

// 自定义错误信息
func (changeTaskStatus ChangeTaskStatus) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"task_id.required": "任务id不能为空",
		"status.required":  "状态不能为空",
		"status.oneof":     "状态值只能是1或者2",
	}
}

type SaveTask struct {
	TaskID     int32  `form:"task_id" json:"task_id"`
	Name       string `form:"name" json:"name" binding:"required"`
	Spec       string `form:"spec" json:"spec" binding:"required"`
	Command    string `form:"command" json:"command" binding:"required"`
	ProcessNum int32  `form:"process_num" json:"process_num" binding:"required"`
	Status     int32  `form:"status" json:"status" binding:"required,oneof=1 2"`
}

func (saveTask SaveTask) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"name.required":        "任务名称不能为空",
		"spec.required":        "任务表达式不能为空",
		"command.required":     "命令不能为空",
		"process_num.required": "进程数不能为空",
		"status.required":      "状态不能为空",
		"status.oneof":         "状态值只能是1或者2",
	}
}

type TaskLogList struct {
	TaskID    int32 `form:"task_id" json:"task_id"`
	Page      int   `form:"page" json:"page" binding:"required"`
	PageSize  int   `form:"page_size" json:"page_size" binding:"required"`
	StartTime struct {
		StartTime string `form:"start_time" json:"start_time"`
		EndTime   string `form:"end_time" json:"end_time"`
	} `form:"start_time" json:"start_time"`
	EndTime struct {
		StartTime string `form:"start_time" json:"start_time"`
		EndTime   string `form:"end_time" json:"end_time"`
	} `form:"end_time" json:"end_time"`
}

func (taskLogList TaskLogList) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		//"task_id.required":   "任务id不能为空",
		"page.required":      "分页不能为空",
		"page_size.required": "分页条数不能为空",
	}
}
