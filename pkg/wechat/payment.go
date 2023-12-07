package wechat

import (
	"billiards/pkg/log"
	"billiards/pkg/tool"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"go.uber.org/zap"
	"path"
	"runtime"
)

// 微信支付相关
// 官方文档：https://github.com/wechatpay-apiv3/wechatpay-go

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
	_, filePath, _, _ := runtime.Caller(0)
	dirPath := path.Dir(filePath)
	certFilePath := dirPath + "/../../cert/apiclient_key.pem"

	p = &Payment{
		appId:               "wx6fab07b9528faf96",
		mchId:               "1660752841",
		notifyUrl:           "https://billiards.wosta.cn/e/wechat/pay-notify",
		mchApiV3Key:         "YLcdCByxLQdCQXNzd3cQ8J8H7dD9P4CU",
		mchCertSerialNumber: "27234D02B606B0557E6361C44D91788FE5FC454A",
		mchCertKeyFile:      certFilePath,
	}

	return
}

// 订单查询
// https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_5_2.shtml
func (p *Payment) GetPaymentOrderDetail(orderNum string) {

}

// 获取退款详情
// https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_5_10.shtml
func (p *Payment) GetRefundDetail(refundNum string) {
	ctx := context.Background()
	client := p.getClient()

	svc := refunddomestic.RefundsApiService{Client: client}
	resp, result, err := svc.QueryByOutRefundNo(ctx,
		refunddomestic.QueryByOutRefundNoRequest{
			OutRefundNo: core.String(refundNum),
		},
	)

	tool.Dump(resp)
	tool.Dump(result)
	if err != nil {
		// 处理错误
		fmt.Printf("call QueryByOutRefundNo err:%s", err)
	}
}

// 发起退款并返回退款结果
// https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_5_9.shtml
func (p *Payment) Refund(amount, total int32, transactionId, outTradeNo, outRefundNo, reason string) (
	resp *refunddomestic.Refund, err error) {

	fmt.Println(amount, total, transactionId, outTradeNo, outRefundNo, reason)
	client := p.getClient()
	ctx := context.Background()

	fmt.Println(12312)

	svc := refunddomestic.RefundsApiService{Client: client}
	resp, result, err := svc.Create(ctx,
		refunddomestic.CreateRequest{
			TransactionId: core.String(transactionId),
			OutTradeNo:    core.String(outTradeNo),
			OutRefundNo:   core.String(outRefundNo),
			Reason:        core.String(reason),
			NotifyUrl:     core.String("https://weixin.qq.com"),
			FundsAccount:  refunddomestic.REQFUNDSACCOUNT_AVAILABLE.Ptr(),
			Amount: &refunddomestic.AmountReq{
				Currency: core.String("CNY"),
				Refund:   core.Int64(int64(amount)),
				Total:    core.Int64(int64(total)),
			},
		},
	)

	tool.Dump(result)
	tool.Dump(resp)
	if err != nil {
		fmt.Println(err.Error())
		log.GetLogger().Error("refund_error", zap.String("msg", err.Error()))
		return
	}

	return
}

// 生成小程序预支付订单（用于前端调起微信支付）
func (p *Payment) GetPrepayBill(openId, description, outTradeNo string, amount int32) (
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
			Total: core.Int64(int64(amount)),
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

func (p *Payment) GetPayResult(c *gin.Context) (
	transaction *payments.Transaction, request *notify.Request, err error) {

	// 1. 使用 `RegisterDownloaderWithPrivateKey` 注册下载器
	mchPrivateKey, _ := utils.LoadPrivateKeyWithPath(p.mchCertKeyFile)
	err = downloader.MgrInstance().
		RegisterDownloaderWithPrivateKey(c, mchPrivateKey, p.mchCertSerialNumber, p.mchId, p.mchApiV3Key)
	if err != nil {
		log.GetLogger().Error("payment_error",
			zap.String("msg", "下载证书失败"),
			zap.String("err", err.Error()))
		err = errors.New("下载证书失败" + err.Error())
		return
	}

	// 2. 获取商户号对应的微信支付平台证书访问器
	certificateVisitor := downloader.MgrInstance().GetCertificateVisitor(p.mchId)

	// 3. 使用证书访问器初始化 `notify.Handler`
	handler, err := notify.NewRSANotifyHandler(p.mchApiV3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))
	if err != nil {
		log.GetLogger().Error("payment_error",
			zap.String("msg", "初始化handler失败"),
			zap.String("err", err.Error()))
	}

	transaction = new(payments.Transaction)

	// 4. 解密数据
	// 如果验签未通过，或者解密失败
	if request, err = handler.ParseNotifyRequest(context.Background(), c.Request, transaction); err != nil {
		log.GetLogger().Error("payment_error",
			zap.String("msg", "数据解密失败"),
			zap.String("err", err.Error()))
		err = errors.New("数据解密失败" + err.Error())
		return
	}

	log.GetLogger().Info("payment_response",
		zap.Any("request", request),
		zap.Any("transaction", transaction))

	return
}
