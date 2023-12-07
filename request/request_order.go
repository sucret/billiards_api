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
