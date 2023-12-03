package model

const (
	OrderStatusDefault      = iota + 1 // 初始状态，待支付
	OrderStatusPaySuccess              // 支付完成（进行中）
	OrderStatusAutoCancel              // 超时未支付系统自动取消
	OrderStatusManualCancel            // 手动取消
	OrderStatusRefund                  // 退款完成
	OrderStatusFinised                 // 订单完成
)
