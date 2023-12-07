// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameCoupon = "coupon"

// Coupon mapped from table <coupon>
type Coupon struct {
	CouponID  int32  `gorm:"column:coupon_id;type:int(11);primaryKey;autoIncrement:true" json:"coupon_id"`
	Name      string `gorm:"column:name;type:varchar(32);not null" json:"name"`
	Price     int32  `gorm:"column:price;type:int(11);not null" json:"price"` // 价格
	ShopID    int32  `gorm:"column:shop_id;type:int(11);not null" json:"shop_id"`
	Type      int    `gorm:"column:type;type:tinyint(4);not null" json:"type"`             // 优惠券类型，1｜抵时卡
	Duration  int32  `gorm:"column:duration;type:int(11);not null" json:"duration"`        // 抵时长，分钟
	TableType int    `gorm:"column:table_type;type:tinyint(4);not null" json:"table_type"` // 类型，1｜通用，2｜球桌，3｜棋牌桌
}

// TableName Coupon's table name
func (*Coupon) TableName() string {
	return TableNameCoupon
}