// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameCouponOrder = "coupon_order"

// CouponOrder mapped from table <coupon_order>
type CouponOrder struct {
	CouponOrderID int32 `gorm:"column:coupon_order_id;type:int(11);primaryKey;autoIncrement:true" json:"coupon_order_id"`
	OrderID       int32 `gorm:"column:order_id;type:int(11);not null" json:"order_id"`
	CouponID      int32 `gorm:"column:coupon_id;type:int(11);not null" json:"coupon_id"`
	PayAmount     int32 `gorm:"column:pay_amount;type:int(11)" json:"pay_amount"` // 订单支付金额
	Status        int   `gorm:"column:status;type:tinyint(4);not null" json:"status"`
	UserID        int32 `gorm:"column:user_id;type:int(11);not null" json:"user_id"`
	CreatedAt     Time  `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     Time  `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// TableName CouponOrder's table name
func (*CouponOrder) TableName() string {
	return TableNameCouponOrder
}
