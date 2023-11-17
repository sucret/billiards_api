package request

type SaveTable struct {
	Name    string `form:"name" json:"name" binding:"required"`
	TableID int32  `form:"table_id" json:"table_id"`
	ShopID  int32  `form:"shop_id" json:"shop_id" binding:"required"`
}

// 自定义错误信息
func (saveShop SaveTable) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"name.required":    "球桌名称不能为空",
		"shop_id.required": "店铺id不能为空",
	}
}
