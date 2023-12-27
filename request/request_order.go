package request

type OrderCreate struct {
	TableID            int32 `form:"table_id" json:"table_id"`
	CouponID           int32 `json:"coupon_id"`
	RechargeAmountType int   `json:"recharge_amount_type"`
	IsRecharge         bool  `json:"is_recharge"`
	UserCouponID       int32 `json:"user_coupon_id"`
}

func (OrderCreate) GetMessages() ValidatorMessages {
	return ValidatorMessages{}
}

type OrderList struct {
	Page     int `form:"page" json:"page" binding:"required"`
	PageSize int `form:"page_size" json:"page_size" binding:"required"`
	ShopId   int `form:"shop_id" json:"shop_id"`
	UserId   int `form:"user_id" json:"user_id"`
}

func (OrderList) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"page.required":      "页码不能为空",
		"page_size.required": "分页条数不能为空",
	}
}
