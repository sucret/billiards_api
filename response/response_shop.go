package response

import "billiards/pkg/mysql/model"

type Shop struct {
	model.Shop
	TableNum     int `json:"table_num"`
	TableFreeNum int `json:"table_free_num"`
}
