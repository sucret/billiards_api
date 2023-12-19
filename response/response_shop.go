package response

import "billiards/pkg/mysql/model"

type Shop struct {
	model.Shop
	BilliardsTableNum     int     `json:"billiards_table_num"`
	BilliardsTableFreeNum int     `json:"billiards_table_free_num"`
	BilliardsPrice        int32   `json:"billiards_price"`
	Distance              float64 `json:"distance"`
}

type TableStatus struct {
	TableID int32 `json:"table_id"`
	Status  int   `json:"status"`
}

type ShopStatusResp struct {
	TableStatusList []TableStatus `json:"table_status_list"`
}
