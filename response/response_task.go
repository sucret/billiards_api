package response

import "gin-api/pkg/mysql/model"

type TaskLogResponse struct {
	Total   int64           `json:"total"`
	LogList []model.TaskLog `json:"list"`
}
