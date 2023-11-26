package model

// 日志类型，1｜发起订单，2｜发起支付，3｜支付完成，4｜取消订单，5｜续费，6｜发起关闭订单，7｜关闭订单失败，8｜关闭订单成功
const (
	OrderLogCreated = iota + 1 // 初始状态，待支付
	OrderLogLaunchPay
	OrderLogPaySuccess
	OrderLogCancel
	OrderLogRenew
	OrderLogLaunchTerminate
	OrderLogTerminateFailed
	OrderLogTerminateSuccess
)
