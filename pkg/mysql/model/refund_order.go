package model

// 退款状态，1｜待退款，2｜退款成功，3｜退款关闭，4｜退款处理中，5｜退款异常
// 微信对应枚举值：SUCCESS|退款成功，CLOSED｜退款关闭，PROCESSING｜退款处理中，ABNORMAL｜退款异常
const (
	RefundStatusDefault  = iota + 1 // 待退款
	RefundStatusSuccess             // 退款成功
	RefundStatusClose               // 退款关闭
	RefundStatusHandling            // 退款处理中
	RefundStatusError               // 退款异常
)

var RefundStatusMapping = map[string]int{
	"SUCCESS":    RefundStatusSuccess,
	"CLOSED":     RefundStatusClose,
	"PROCESSING": RefundStatusHandling,
	"ABNORMAL":   RefundStatusError,
}
