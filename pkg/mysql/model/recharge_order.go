package model

const (
	// 订单状态，1｜待支付，2｜支付完成，3｜支付取消
	RechargeOrderStatusDefault = iota + 1
	RechargeOrderStatusPaid
	RechargeOrderStatusCancel
)
