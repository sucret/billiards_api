package response

import "billiards/pkg/mysql/model"

type Shop struct {
	model.Shop
	BilliardsTableNum     int   `json:"billiards_table_num"`
	BilliardsTableFreeNum int   `json:"billiards_table_free_num"`
	BilliardsPrice        int32 `json:"billiards_price"`
}
