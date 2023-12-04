package model

const (
	POTypeTable    = iota + 1 // 开台订单
	POTypeRecharge            // 充值订单
)

// 支付方式
const (
	PMOModeWechat = iota + 1 // 微信支付
	PMOModeWallet            // 钱包余额支付
)

// 支付状态
const (
	PMOStatusDefault = iota + 1 // 默认状态，待支付
	PMOStatusSuccess            // 支付成功
)
