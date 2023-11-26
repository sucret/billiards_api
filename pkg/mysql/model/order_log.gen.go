// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameOrderLog = "order_log"

// OrderLog mapped from table <order_log>
type OrderLog struct {
	OrderLogID int32  `gorm:"column:order_log_id;type:int(11);primaryKey;autoIncrement:true" json:"order_log_id"`
	OrderID    int32  `gorm:"column:order_id;type:int(11);not null" json:"order_id"`
	Type       int32  `gorm:"column:type;type:tinyint(4);not null" json:"type"` // 日志类型，1｜发起订单，2｜发起支付，3｜支付完成，4｜取消订单，5｜续费，6｜发起关闭订单，7｜关闭订单失败，8｜关闭订单成功
	Remark     string `gorm:"column:remark;type:text" json:"remark"`
	CreatedAt  Time   `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName OrderLog's table name
func (*OrderLog) TableName() string {
	return TableNameOrderLog
}
