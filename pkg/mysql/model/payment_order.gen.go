// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNamePaymentOrder = "payment_order"

// PaymentOrder mapped from table <payment_order>
type PaymentOrder struct {
	PaymentOrderID int32   `gorm:"column:payment_order_id;type:int(11);primaryKey;autoIncrement:true" json:"payment_order_id"`
	OrderID        int32   `gorm:"column:order_id;type:int(11);not null" json:"order_id"`
	OrderNum       string  `gorm:"column:order_num;type:char(20);not null" json:"order_num"`
	Amount         float64 `gorm:"column:amount;type:decimal(10,2);not null;default:0.00" json:"amount"`
	CreatedAt      Time    `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      Time    `gorm:"column:updated_at;type:datetime" json:"updated_at"`
	NotifyID       string  `gorm:"column:notify_id;type:char(100);not null" json:"notify_id"` // 微信通知id
	Resource       string  `gorm:"column:resource;type:text" json:"resource"`
}

// TableName PaymentOrder's table name
func (*PaymentOrder) TableName() string {
	return TableNamePaymentOrder
}
