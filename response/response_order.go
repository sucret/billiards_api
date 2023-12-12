package response

import (
	"billiards/pkg/mysql/model"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
)

type OrderDetail struct {
	model.TableOrder
	UsedMinutes   int32 `json:"used_minutes"`
	RemainMinutes int32 `json:"remain_minutes"`
	TotalMinutes  int32 `json:"total_minutes"`
}

type RechargeResult struct {
	Succeed bool  `json:"succeed"`
	Wallet  int32 `json:"wallet"`
}

type PayResult struct {
	Succeed bool  `json:"succeed"`
	Wallet  int32 `json:"wallet"`
}

type PaymentOrderResp struct {
	Order              model.Order
	WxPaymentOrder     model.PaymentOrder
	WalletPaymentOrder model.PaymentOrder
	WxPayResp          *jsapi.PrepayWithRequestPaymentResponse
}

type OrderResp struct {
	Order     model.Order                             `json:"order"`
	WxPayResp *jsapi.PrepayWithRequestPaymentResponse `json:"wx_pay_resp"`
	NeedWxPay bool                                    `json:"need_wx_pay"`
	//Payment
}

type OrderDetailResp struct {
	Order            model.Order       `json:"order"`
	TableOrder       OrderDetail       `json:"table_order"`
	CouponOrder      model.CouponOrder `json:"coupon_order"`
	PaymentOrderList []model.PaymentOrder
}
