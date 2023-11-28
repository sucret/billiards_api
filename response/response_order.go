package response

import (
	"billiards/pkg/mysql/model"
)

type OrderDetail struct {
	model.Order
	UsedMinutes   int32 `json:"used_minutes"`
	RemainMinutes int32 `json:"remain_minutes"`
	TotalMinutes  int32 `json:"total_minutes"`
}
