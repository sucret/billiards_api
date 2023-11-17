package response

import "billiards/pkg/mysql/model"

type TaskLogResponse struct {
	Total   int64           `json:"total"`
	LogList []model.TaskLog `json:"list"`
}
