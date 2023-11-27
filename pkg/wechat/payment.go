package wechat

import (
	"billiards/pkg/log"
	"billiards/pkg/tool"
	"context"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"go.uber.org/zap"
	"net/http"
)

type Payment struct {
	appId               string // 小程序appid
	mchId               string // 商户号
	notifyUrl           string // 支付回调地址
	mchApiV3Key         string // 商户APIv3密钥
	mchCertSerialNumber string // 商户证书序列号
	mchCertKeyFile      string // 证书文件
}

func NewPayment() (p *Payment) {
	// TODO 参数后期加入配置文件里边
	p = &Payment{
		appId:               "wx6fab07b9528faf96",
		mchId:               "1660752841",
		notifyUrl:           "http://http://billiards.wosta.cn/e/wechat/pay-notify",
		mchApiV3Key:         "YLcdCByxLQdCQXNzd3cQ8J8H7dD9P4CU",
		mchCertSerialNumber: "27234D02B606B0557E6361C44D91788FE5FC454A",
		mchCertKeyFile:      "cert/apiclient_key.pem",
	}

	return
}

// 生成小程序预支付订单（用于前端调起微信支付）
func (p *Payment) GetPrepayBill(openId, description, outTradeNo string, amount int64) (
	res *jsapi.PrepayWithRequestPaymentResponse, err error) {
	client := p.getClient()
	svc := jsapi.JsapiApiService{Client: client}

	ctx := context.Background()

	param := jsapi.PrepayRequest{
		Appid:       core.String(p.appId),
		Mchid:       core.String(p.mchId),
		Description: core.String(description),
		OutTradeNo:  core.String(outTradeNo),
		Attach:      core.String(""),
		NotifyUrl:   core.String(p.notifyUrl),
		Amount: &jsapi.Amount{
			Total: core.Int64(amount),
		},
		Payer: &jsapi.Payer{
			Openid: core.String(openId),
		},
	}

	if res, _, err = svc.PrepayWithRequestPayment(ctx, param); err != nil {
		log.GetLogger().Error("payment", zap.String("msg", err.Error()), zap.Any("param", param))
	}

	return
}

func (p *Payment) getClient() (client *core.Client) {
	mchPrivateKey, err := utils.LoadPrivateKeyWithPath(p.mchCertKeyFile)
	if err != nil {
		log.GetLogger().Error("payment", zap.String("msg", err.Error()))
	}

	ctx := context.Background()
	// 使用商户私钥等初始化 client，并使它具有自动定时获取微信支付平台证书的能力
	opts := []core.ClientOption{
		option.WithWechatPayAutoAuthCipher(p.mchId, p.mchCertSerialNumber, mchPrivateKey, p.mchApiV3Key),
	}
	client, err = core.NewClient(ctx, opts...)
	if err != nil {
		log.GetLogger().Error("payment", zap.String("msg", err.Error()), zap.Any("param", p))
	}

	return
}

// 官方文档：https://pay.weixin.qq.com/docs/partner/apis/partner-jsapi-payment/payment-notice.html
func (p *Payment) GetPayResult(r *http.Request) (interface{}, int32) {
	fmt.Println("微信支付回调")

	handler, err := notify.NewRSANotifyHandler(
		p.mchApiV3Key, verifiers.NewSHA256WithRSAVerifier(core.NewCertificateMapWithList(nil)),
	)

	content := make(map[string]interface{})
	request, err := handler.ParseNotifyRequest(context.Background(), r, content)
	tool.Dump(err)

	fmt.Println(request)
	tool.Dump(request)
	fmt.Println(content)
	tool.Dump(content)
	return "success", 0
}
