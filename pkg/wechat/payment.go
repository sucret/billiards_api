package wechat

import (
	"billiards/pkg/log"
	"billiards/pkg/tool"
	"context"
	"crypto/x509"
	"fmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/core/option"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"github.com/wechatpay-apiv3/wechatpay-go/utils"
	"go.uber.org/zap"
	"net/http"
	"time"
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
		notifyUrl:           "https://billiards.wosta.cn/e/wechat/pay-notify",
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

type contentType struct {
	Mchid           *string    `json:"mchid"`
	Appid           *string    `json:"appid"`
	CreateTime      *time.Time `json:"create_time"`
	OutContractCode *string    `json:"out_contract_code"`
}

// 官方文档：https://pay.weixin.qq.com/docs/partner/apis/partner-jsapi-payment/payment-notice.html
func (p *Payment) GetPayResult(r *http.Request) (interface{}, int32) {

	mchAPIv3Key := p.mchApiV3Key

	wechatPayCert, err := utils.LoadCertificateWithPath("cert/apiclient_cert.pem")
	//wechatPayCert, err := utils.LoadCertificate("<your wechat pay certificate>")
	// 2. 使用本地管理的微信支付平台证书获取微信支付平台证书访问器
	certificateVisitor := core.NewCertificateMapWithList([]*x509.Certificate{wechatPayCert})
	// 3. 使用apiv3 key、证书访问器初始化 `notify.Handler`
	handler := notify.NewNotifyHandler(mchAPIv3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))

	transaction := new(payments.Transaction)

	notifyReq, err := handler.ParseNotifyRequest(context.Background(), r, transaction)
	// 如果验签未通过，或者解密失败
	if err != nil {
		fmt.Println(333, err)
		//return
	}

	fmt.Println(notifyReq)
	fmt.Println(transaction)
	tool.Dump(notifyReq)
	tool.Dump(transaction)

	//fmt.Println("微信支付回调")
	//
	//cert, err := utils.LoadCertificateWithPath("cert/apiclient_cert.pem")
	////cert, err := utils.LoadCertificate("-----BEGIN CERTIFICATE-----\nMIIEITCCAwmgAwIBAgIUJyNNArYGsFV+Y2HETZF4j+X8RUowDQYJKoZIhvcNAQEL\nBQAwXjELMAkGA1UEBhMCQ04xEzARBgNVBAoTClRlbnBheS5jb20xHTAbBgNVBAsT\nFFRlbnBheS5jb20gQ0EgQ2VudGVyMRswGQYDVQQDExJUZW5wYXkuY29tIFJvb3Qg\nQ0EwHhcNMjMxMTI3MDI1NDMyWhcNMjgxMTI1MDI1NDMyWjB7MRMwEQYDVQQDDAox\nNjYwNzUyODQxMRswGQYDVQQKDBLlvq7kv6HllYbmiLfns7vnu58xJzAlBgNVBAsM\nHumZleilv+S6rOagvOagh+ivhuaciemZkOWFrOWPuDELMAkGA1UEBhMCQ04xETAP\nBgNVBAcMCFNoZW5aaGVuMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA\nvIVBHBDNRZVO7C6d5TUjZFNIn8eGO5ztRJ1yyg1566bEjN8iheU7SxOYkXr4beSx\nEDyYom2oSLg63+AHfPSdHNqyC/id4qlJXySWZQgUK5yn1Vp5xHL5i0yaHEcBsjEm\njDOxuf4elD1m2xu6djsq0x7UzGvphyNYrUnxS33STnrhDFn/zWl6lF2LOYpR5vQ+\ne2Pe33kozM55bx4CgGEms56Z4QSY69jy6b8yww1g2y/QHqhQ1aUwzeqcPd5S+t/9\nLaMyGAn5Ty9gwDsXC0RFFJcBL7tpTO2lcZJegPwAZplPmPINPsqep5z8W3GaGl2H\na2Br/TStYPcUrKgNtmC1XwIDAQABo4G5MIG2MAkGA1UdEwQCMAAwCwYDVR0PBAQD\nAgP4MIGbBgNVHR8EgZMwgZAwgY2ggYqggYeGgYRodHRwOi8vZXZjYS5pdHJ1cy5j\nb20uY24vcHVibGljL2l0cnVzY3JsP0NBPTFCRDQyMjBFNTBEQkMwNEIwNkFEMzk3\nNTQ5ODQ2QzAxQzNFOEVCRDImc2c9SEFDQzQ3MUI2NTQyMkUxMkIyN0E5RDMzQTg3\nQUQxQ0RGNTkyNkUxNDAzNzEwDQYJKoZIhvcNAQELBQADggEBAA7j2iUxdwCjMMFB\nkIdBQFbqwXJ6yftKrP+0AefVd44XwvS7+zdIdAnRDVjs5eWVv8+W0VhE6+jnu06K\nIoC0MvuCAqz1rmFkVpZr6Qyw+jjKM4u+Bf3lWSy2tWfVnZNEn64BJOtiXNrAxhn+\nSsTcXYW9pBUj+VeKMFnsGgkGY7zF+a4f75GFhHRySEnEnJw4v18uIIgrAIU1qZhL\nmvGLenTaB/Ppd4eW/Usp2JQ3MsLPtLUe7/ZqYMkPrN0RVyPHMShpcF6nynGkZkAe\nOU1P5lVGnAnPGwe22JtTspxpQHSyAQzQU0o5COu56qLHudT/n9r3XFzcrjL6zkWR\nIDGtBKI=\n-----END CERTIFICATE-----")
	//if err != nil {
	//	fmt.Println("err", err.Error())
	//}
	//
	////handler, err := notify.NewRSANotifyHandler(
	////	p.mchCertKeyFile, verifiers.NewSHA256WithRSAVerifier(core.NewCertificateMapWithList([]*x509.Certificate{cert})),
	////)
	//
	//handler, _ := notify.NewRSANotifyHandler(p.mchApiV3Key, verifiers.NewSHA256WithRSAVerifier(core.NewCertificateMapWithList(nil)))
	//
	//fmt.Println(1234, err.Error())
	//
	//fmt.Println(handler)

	//content := contentType{}
	//
	//notifyReq, err := handler.ParseNotifyRequest(context.Background(), r, content)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(err, 3224)
	//tool.Dump(notifyReq)
	//wechatPayCert, err := utils.LoadCertificate("-----BEGIN CERTIFICATE-----\nMIIEITCCAwmgAwIBAgIUJyNNArYGsFV+Y2HETZF4j+X8RUowDQYJKoZIhvcNAQEL\nBQAwXjELMAkGA1UEBhMCQ04xEzARBgNVBAoTClRlbnBheS5jb20xHTAbBgNVBAsT\nFFRlbnBheS5jb20gQ0EgQ2VudGVyMRswGQYDVQQDExJUZW5wYXkuY29tIFJvb3Qg\nQ0EwHhcNMjMxMTI3MDI1NDMyWhcNMjgxMTI1MDI1NDMyWjB7MRMwEQYDVQQDDAox\nNjYwNzUyODQxMRswGQYDVQQKDBLlvq7kv6HllYbmiLfns7vnu58xJzAlBgNVBAsM\nHumZleilv+S6rOagvOagh+ivhuaciemZkOWFrOWPuDELMAkGA1UEBhMCQ04xETAP\nBgNVBAcMCFNoZW5aaGVuMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA\nvIVBHBDNRZVO7C6d5TUjZFNIn8eGO5ztRJ1yyg1566bEjN8iheU7SxOYkXr4beSx\nEDyYom2oSLg63+AHfPSdHNqyC/id4qlJXySWZQgUK5yn1Vp5xHL5i0yaHEcBsjEm\njDOxuf4elD1m2xu6djsq0x7UzGvphyNYrUnxS33STnrhDFn/zWl6lF2LOYpR5vQ+\ne2Pe33kozM55bx4CgGEms56Z4QSY69jy6b8yww1g2y/QHqhQ1aUwzeqcPd5S+t/9\nLaMyGAn5Ty9gwDsXC0RFFJcBL7tpTO2lcZJegPwAZplPmPINPsqep5z8W3GaGl2H\na2Br/TStYPcUrKgNtmC1XwIDAQABo4G5MIG2MAkGA1UdEwQCMAAwCwYDVR0PBAQD\nAgP4MIGbBgNVHR8EgZMwgZAwgY2ggYqggYeGgYRodHRwOi8vZXZjYS5pdHJ1cy5j\nb20uY24vcHVibGljL2l0cnVzY3JsP0NBPTFCRDQyMjBFNTBEQkMwNEIwNkFEMzk3\nNTQ5ODQ2QzAxQzNFOEVCRDImc2c9SEFDQzQ3MUI2NTQyMkUxMkIyN0E5RDMzQTg3\nQUQxQ0RGNTkyNkUxNDAzNzEwDQYJKoZIhvcNAQELBQADggEBAA7j2iUxdwCjMMFB\nkIdBQFbqwXJ6yftKrP+0AefVd44XwvS7+zdIdAnRDVjs5eWVv8+W0VhE6+jnu06K\nIoC0MvuCAqz1rmFkVpZr6Qyw+jjKM4u+Bf3lWSy2tWfVnZNEn64BJOtiXNrAxhn+\nSsTcXYW9pBUj+VeKMFnsGgkGY7zF+a4f75GFhHRySEnEnJw4v18uIIgrAIU1qZhL\nmvGLenTaB/Ppd4eW/Usp2JQ3MsLPtLUe7/ZqYMkPrN0RVyPHMShpcF6nynGkZkAe\nOU1P5lVGnAnPGwe22JtTspxpQHSyAQzQU0o5COu56qLHudT/n9r3XFzcrjL6zkWR\nIDGtBKI=\n-----END CERTIFICATE-----")
	//if err != nil {
	//	fmt.Println(err)
	//
	//}
	//fmt.Println(wechatPayCert, 123)
	//
	//certificateVisitor := core.NewCertificateMapWithList([]*x509.Certificate{wechatPayCert})
	//handler := notify.NewNotifyHandler(p.mchApiV3Key, verifiers.NewSHA256WithRSAVerifier(certificateVisitor))

	//handler, err := notify.NewRSANotifyHandler(
	//	p.mchApiV3Key, verifiers.NewSHA256WithRSAVerifier(core.NewCertificateMapWithList(nil)),
	//)
	//
	//content := make(map[string]interface{})
	////content := new(payments.Transaction)
	//
	//request, err := handler.ParseNotifyRequest(context.Background(), r, content)
	//tool.Dump(err)
	//
	//fmt.Println(request)
	//tool.Dump(request)
	//fmt.Println(content)
	//tool.Dump(content)
	return "success", 0
}
