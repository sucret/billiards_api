package service

import (
	"billiards/pkg/log"
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	redis_ "billiards/pkg/redis"
	"billiards/pkg/tool"
	"billiards/pkg/wechat"
	"billiards/response"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type paymentService struct {
	db    *gorm.DB
	redis *redis.Client
}

var PaymentService = &paymentService{
	db:    mysql.GetDB(),
	redis: redis_.GetRedis(),
}

func (p *paymentService) refund(amount int32, orderList *[]model.PaymentOrder) {
	for _, v := range *orderList {
		// 付款金额如果大于退款金额，则退掉就结束
		// 付款金额如果
		payAmount := v.Amount

		if amount <= 0 {
			break
		}

		refundAmount := amount

		// 如果支付金额小于退款金额，则这一单整单退，然后剩余的钱下一单给退
		if v.Amount < amount {
			refundAmount = payAmount
			amount = amount - refundAmount
		}

		// 退款
		if v.PayMode == model.PMOModeWallet {
			// 退还到钱包
			_, err := p.RefundWalletPay(v, refundAmount)
			if err != nil {
				return
			}
		} else {
			// 微信支付的原路退回
			_, err := p.RefundWxPay(v, refundAmount)
			if err != nil {
				return
			}
		}
	}

	if amount > 0 {
		log.GetLogger().Error("refund_error", zap.String("msg", "退款金额超出支付金额"))
	}
}

// 生成付款单
// 可以是普通下单，也可以是续费支付
// 根据能否使用余额以及用户余额来生成付款单
// 如果需要用微信支付，则生成微信预支付参数
func (p *paymentService) createOrder(tx *gorm.DB, userId, orderId, payAmount int32, canUseWallet, isRenewal bool) (
	resp response.PaymentOrderResp, err error) {

	// 1、计算微信支付和余额支付各自的金额
	var wxPayAmount int32     // 微信支付金额
	var walletPayAmount int32 // 钱包支付金额

	user, err := UserService.GetByUserId(userId)
	if err != nil {
		return response.PaymentOrderResp{}, err
	}
	if canUseWallet {

		if user.Wallet > payAmount {
			// 钱包大于支付金额，则全额用钱包支付
			walletPayAmount = payAmount
		} else if user.Wallet > 0 && user.Wallet < payAmount {
			// 否则钱包支付部分，微信支付剩余部分
			walletPayAmount = user.Wallet
			wxPayAmount = payAmount - walletPayAmount
		}
	} else {
		wxPayAmount = payAmount
	}

	// 2、有余额付款则生成余额付款单
	if walletPayAmount > 0 {
		resp.WalletPaymentOrder, err = p.makePaymentOrder(tx, walletPayAmount, userId, orderId, model.PMOModeWallet)
		if err != nil {
			return response.PaymentOrderResp{}, err
		}
	}

	// 3、有微信支付则生成微信付款单
	if wxPayAmount > 0 {
		resp.WxPaymentOrder, err = p.makePaymentOrder(tx, wxPayAmount, userId, orderId, model.PMOModeWechat)
		if err != nil {
			return response.PaymentOrderResp{}, err
		}

		// 微信支付附加信息（支付成功之后会原样回调给回调接口，用于处理回调部分的逻辑）
		attach := &wechat.Attach{OrderId: orderId, IsRenewal: isRenewal}
		// 生成微信支付的参数
		resp.WxPayResp, err = wechat.NewPayment().GetPrepayBill(
			user.OpenID, "description", resp.WxPaymentOrder.PaymentOrderNo, wxPayAmount, attach)
		if err != nil {
			return response.PaymentOrderResp{}, err
		}
	}

	// 4、没有微信支付的话，那就直接把订单改成已支付
	if wxPayAmount == 0 && walletPayAmount >= 0 {
		if isRenewal {
			err = p.doOrderSuccess(tx, &resp.WalletPaymentOrder)
		} else {
			err = OrderService.PaySuccess(tx, orderId)
			if err != nil {
				return response.PaymentOrderResp{}, err
			}
		}
	}

	return
}

// 给指定付款单退款
func (p *paymentService) RefundWxPay(paymentOrder model.PaymentOrder, amount int32) (
	refundOrder model.RefundOrder, err error) {
	// 1、创建退款单
	//Amount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", amount/100), 64)
	//totalAmount, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", paymentOrder.Amount), 64)

	// 2、发起退款
	// 3、回写退款单信息
	// 4、返回退款结果

	refundOrder = model.RefundOrder{
		OrderID:        paymentOrder.OrderID,
		PaymentOrderID: paymentOrder.PaymentOrderID,
		Amount:         amount,
		//Status:         int32(model.RefundStatusMapping[string(*resp.Status)]),
		RefundNum: tool.GenerateOrderNum(),
		//WxRefundID:     *resp.RefundId,
	}

	resp, err := wechat.NewPayment().Refund(amount, paymentOrder.Amount, paymentOrder.TransactionID,
		paymentOrder.PaymentOrderNo, refundOrder.RefundNum, "退款")

	refundOrder.Status = model.RefundStatusMapping[string(*resp.Status)]
	refundOrder.WxRefundID = *resp.RefundId

	err = p.db.Create(&refundOrder).Error

	return
}

// 给钱包付款单退款
func (p *paymentService) RefundWalletPay(paymentOrder model.PaymentOrder, amount int32) (
	refundOrder model.RefundOrder, err error) {
	// 获取用户信息
	user := model.User{}
	p.db.Where("user_id = ?", paymentOrder.UserID).First(&user)

	// 更新用户余额
	user.Wallet = user.Wallet + amount
	p.db.Save(&user)

	refundOrder = model.RefundOrder{
		OrderID:        paymentOrder.OrderID,
		PaymentOrderID: paymentOrder.PaymentOrderID,
		Amount:         amount,
		Status:         model.RefundStatusMapping["SUCCESS"],
		//Status:         int32(model.RefundStatusMapping[string(*resp.Status)]),
		RefundNum: tool.GenerateOrderNum(),
		//WxRefundID:     *resp.RefundId,
	}

	// 写入记录
	err = p.db.Create(&refundOrder).Error

	return
}

// 给指定付款单退款
func (p *paymentService) RefundOrder(order model.TableOrder, amount int32) {
	//for _, v := range order.PaymentOrderList {
	//	// 付款金额如果大于退款金额，则退掉就结束
	//	// 付款金额如果
	//	payAmount := v.Amount
	//
	//	if amount <= 0 {
	//		break
	//	}
	//
	//	refundAmount := amount
	//
	//	// 如果支付金额小于退款金额，则这一单整单退，然后剩余的钱下一单给退
	//	if payAmount < refundAmount {
	//		refundAmount = payAmount
	//		amount = amount - refundAmount
	//	}
	//
	//	// 退款
	//	if v.PayMode == model.PMOModeWallet {
	//		// 退还到钱包
	//		_, err := p.RefundWalletPay(v, order.UserID, refundAmount)
	//		if err != nil {
	//			return
	//		}
	//	} else {
	//		// 微信支付的原路退回
	//		_, err := p.RefundWxPay(v, refundAmount)
	//		if err != nil {
	//			return
	//		}
	//	}
	//}
	//
	//if amount > 0 {
	//	log.GetLogger().Error("refund_error", zap.String("msg", "退款金额超出支付金额"))
	//}

	return
}

// 付款单支付成功逻辑
func (p *paymentService) doOrderSuccess(db *gorm.DB, paymentOrder *model.PaymentOrder) (err error) {
	// 更新付款单状态
	paymentOrder.Status = model.PMOStatusSuccess
	if err = db.Save(&paymentOrder).Error; err != nil {
		log.GetLogger().Error("pay_error", zap.String("msg", "更新支付订单失败："+err.Error()))
		return
	}

	// 如果是钱包支付，这里需要去扣掉用户余额
	if paymentOrder.PayMode == model.PMOModeWallet {
		err = UserService.ChangeWallet(db, paymentOrder.UserID, paymentOrder.Amount*-1)
		if err != nil {
			return err
		}
	}

	return
}

// 微信支付回调
func (p *paymentService) WechatPayNotify(c *gin.Context) (err error) {
	res, request, err := wechat.NewPayment().GetPayResult(c)
	if err != nil {
		return
	}

	tx := p.db.Begin()
	defer tx.Rollback()

	// 根据微信返回的付款单号查询对应的付款单
	paymentOrder := model.PaymentOrder{}
	if err = tx.Where("payment_order_no = ?", res.OutTradeNo).First(&paymentOrder).Error; err != nil {
		log.GetLogger().Error("pay_error", zap.String("msg", "查询支付订单失败："+err.Error()))
		return
	}

	if *res.TradeState != "SUCCESS" {
		log.GetLogger().Error("pay_error", zap.String("msg", "订单支付失败"+*res.TradeState))
		return
	}

	// 更新付款单信息
	paymentOrder.NotifyID = request.ID
	paymentOrder.Resource = ""
	paymentOrder.BankType = *res.BankType
	paymentOrder.TransactionID = *res.TransactionId
	paymentOrder.TradeState = *res.TradeState
	if err = tx.Save(&paymentOrder).Error; err != nil {
		log.GetLogger().Error("pay_error", zap.String("msg", "更新付款单信息失败："+err.Error()))
		return
	}

	attach := wechat.Attach{}
	err = json.Unmarshal([]byte(*res.Attach), &attach)
	if err != nil {
		log.GetLogger().Error("pay_error", zap.String("msg", "attach转struct失败，"+err.Error()))
	}

	// 调用订单支付成功的方法
	if attach.IsRenewal {
		// 如果是续费，则直接执行付款单成功的逻辑
		err = p.doOrderSuccess(tx, &paymentOrder)
		if err != nil {
			return err
		}
	} else {
		// 正常业务单据走业务单据的支付成功逻辑
		err = OrderService.PaySuccess(tx, paymentOrder.OrderID)
		if err != nil {
			log.GetLogger().Error("pay_error", zap.String("msg", "更新订单失败："+err.Error()))
			return err
		}
	}

	tx.Commit()

	return
}

// 创建优惠券订单
func (p *paymentService) MakeCouponOrder(db *gorm.DB, amount int32, orderType int,
	orderId, userId int32, description string) (
	payment *jsapi.PrepayWithRequestPaymentResponse, err error) {

	//user, _ := UserService.GetByUserId(userId)
	//
	//paymentOrder := model.PaymentOrder{
	//	PaymentOrderNo: tool.GenerateOrderNum(),
	//	OrderID:        orderId,
	//	Amount:         amount,
	//	OrderType:      orderType,
	//
	//	PayMode: model.PMOModeWechat,
	//	Status:  model.PMOStatusDefault,
	//}
	//
	//if err = db.Create(&paymentOrder).Error; err != nil {
	//	return
	//}

	//payment, err = wechat.NewPayment().GetPrepayBill(user.OpenID, description, paymentOrder.PaymentOrderNo, amount)
	if err != nil {
		return
	}

	return
}

//创建余额支付订单
func (p *paymentService) makePaymentOrder(db *gorm.DB, amount, userId, orderId int32, payMode int) (
	order model.PaymentOrder, err error) {

	order = model.PaymentOrder{
		PaymentOrderNo: tool.GenerateOrderNum(),
		OrderID:        orderId,
		Amount:         amount,
		PayMode:        payMode,
		Status:         model.PMOStatusDefault,
		UserID:         userId,
		CreatedAt:      model.Time(time.Now()),
	}

	if err = db.Create(&order).Error; err != nil {
		return
	}

	return
}
