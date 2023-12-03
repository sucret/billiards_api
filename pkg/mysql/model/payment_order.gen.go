// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNamePaymentOrder = "payment_order"

// PaymentOrder mapped from table <payment_order>
type PaymentOrder struct {
	PaymentOrderID int32  `gorm:"column:payment_order_id;type:int(11);primaryKey;autoIncrement:true" json:"payment_order_id"`
	OrderID        int32  `gorm:"column:order_id;type:int(11);not null" json:"order_id"`
	OrderType      int32  `gorm:"column:order_type;type:tinyint(4);not null" json:"order_type"`           // 订单类型，1｜开台订单支付，2｜充值支付
	PaymentOrderNo string `gorm:"column:payment_order_no;type:char(20);not null" json:"payment_order_no"` // 订单号
	Amount         int32  `gorm:"column:amount;type:int(11);not null" json:"amount"`
	NotifyID       string `gorm:"column:notify_id;type:char(100);not null" json:"notify_id"` // 微信通知id
	Resource       string `gorm:"column:resource;type:text" json:"resource"`
	BankType       string `gorm:"column:bank_type;type:char(32);not null" json:"bank_type"`           // 银行
	TransactionID  string `gorm:"column:transaction_id;type:char(32);not null" json:"transaction_id"` // 微信支付系统订单号
	/*
		交易状态，枚举值：
		SUCCESS：支付成功
		REFUND：转入退款
		NOTPAY：未支付
		CLOSED：已关闭
		REVOKED：已撤销（付款码支付）
		USERPAYING：用户支付中（付款码支付）
		PAYERROR：支付失败(其他原因，如银行返回失败)
	*/
	TradeState string `gorm:"column:trade_state;type:char(12);not null" json:"trade_state"`
	CreatedAt  Time   `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  Time   `gorm:"column:updated_at;type:datetime" json:"updated_at"`
}

// TableName PaymentOrder's table name
func (*PaymentOrder) TableName() string {
	return TableNamePaymentOrder
}
