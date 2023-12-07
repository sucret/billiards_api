package request

type SaveCoupon struct {
	CouponID  int32  `json:"coupon_id"`
	Name      string `json:"name"`
	Price     int32  `json:"price"`
	ShopID    int32  `json:"shop_id"`
	Duration  int32  `json:"duration"`
	Type      int    `json:"type"`
	TableType int    `json:"table_type"`
}
