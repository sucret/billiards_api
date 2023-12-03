package response

import (
	"billiards/pkg/mysql/model"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
)

type TableOrderPrePayParam struct {
	JsApi *jsapi.PrepayWithRequestPaymentResponse `json:"js_api"`
	Order *model.TableOrder                       `json:"order"`
}

type RechargeOrderPrePayParam struct {
	JsApi *jsapi.PrepayWithRequestPaymentResponse `json:"js_api"`
	Order *model.RechargeOrder                    `json:"order"`
}
