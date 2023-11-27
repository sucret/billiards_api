package service

import (
	"billiards/pkg/log"
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	redis_ "billiards/pkg/redis"
	"billiards/pkg/wechat"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type paymentService struct {
	db    *gorm.DB
	redis *redis.Client
}

var PaymentService = &paymentService{
	db:    mysql.GetDB(),
	redis: redis_.GetRedis(),
}

func (p *paymentService) PayNotify(c *gin.Context) (err error) {
	res, request, err := wechat.NewPayment().GetPayResult(c)
	if err != nil {
		return
	}

	paymentOrder := model.PaymentOrder{}
	if err = p.db.Where("order_num = ?", res.OutTradeNo).First(&paymentOrder).Error; err != nil {
		log.GetLogger().Error("pay_error", zap.String("msg", "查询支付订单失败："+err.Error()))
		return
	}

	paymentOrder.NotifyID = request.ID
	paymentOrder.Resource = ""
	paymentOrder.BankType = *res.BankType
	paymentOrder.TransactionID = *res.TransactionId
	paymentOrder.TradeState = *res.TradeState
	fmt.Println(res, request)

	if err = p.db.Save(&paymentOrder).Error; err != nil {
		log.GetLogger().Error("pay_error", zap.String("msg", "更新支付订单失败："+err.Error()))
		return
	}

	_, err = OrderService.PaySuccess(paymentOrder.OrderNum)
	if err != nil {
		log.GetLogger().Error("pay_error", zap.String("msg", "更新订单失败："+err.Error()))
		return err
	}

	return
}

// 创建预支付订单并生成预支付参数
func (p *paymentService) MakePrepayOrder(order model.Order, openID, description, orderNum string, amount float64) (
	payment *jsapi.PrepayWithRequestPaymentResponse, err error) {

	// 写入payment_order表
	paymentOrder := model.PaymentOrder{
		OrderNum: orderNum,
		OrderID:  order.OrderID,
		Amount:   amount,
	}

	if err = p.db.Create(&paymentOrder).Error; err != nil {
		return
	}

	payment, err = wechat.NewPayment().GetPrepayBill(openID, description, orderNum, int64(amount*100))
	if err != nil {
		return
	}

	return
}
