// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameRefundOrder = "refund_order"

// RefundOrder mapped from table <refund_order>
type RefundOrder struct {
	RefundOrderID  int32  `gorm:"column:refund_order_id;type:int(11);primaryKey;autoIncrement:true" json:"refund_order_id"`
	PaymentOrderID int32  `gorm:"column:payment_order_id;type:int(11);not null" json:"payment_order_id"` // 付款单id
	OrderID        int32  `gorm:"column:order_id;type:int(11);not null" json:"order_id"`
	OrderNum       string `gorm:"column:order_num;type:char(32);not null" json:"order_num"`
	Status         int32  `gorm:"column:status;type:tinyint(4);not null" json:"status"`           // 退款状态，1｜待退款，2｜退款成功，3｜退款关闭，4｜退款处理中，5｜退款异常，微信对应枚举值：SUCCESS|退款成功，CLOSED｜退款关闭，PROCESSING｜退款处理中，ABNORMAL｜退款异常
	RefundNum      string `gorm:"column:refund_num;type:char(32);not null" json:"refund_num"`     // 退款单号，不可重复，同意单号多次请求只会退一次
	Amount         int32  `gorm:"column:amount;type:int(11);not null" json:"amount"`              // 退款金额
	WxRefundID     string `gorm:"column:wx_refund_id;type:char(32);not null" json:"wx_refund_id"` // 微信支付退款单号
	CreatedAt      Time   `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      Time   `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// TableName RefundOrder's table name
func (*RefundOrder) TableName() string {
	return TableNameRefundOrder
}