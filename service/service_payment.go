package service

import (
	"billiards/pkg/log"
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	redis_ "billiards/pkg/redis"
	"billiards/pkg/tool"
	"billiards/pkg/wechat"
	"billiards/response"
	"errors"
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

// 生成付款单
// 可以是普通下单，也可以是续费支付
// 根据能否使用余额以及用户余额来生成付款单
// 如果需要用微信支付，则生成微信预支付参数
func (p *paymentService) CreateOrder(tx *gorm.DB, userId, orderId, payAmount int32, canUseWallet bool) (
	resp response.PaymentOrderResp, err error) {
	// 计算微信支付和余额支付各自的金额
	var wxPayAmount int32
	var walletPayAmount int32

	user, err := UserService.GetByUserId(userId)
	if err != nil {
		return response.PaymentOrderResp{}, err
	}

	if canUseWallet {
		if user.Wallet > payAmount {
			walletPayAmount = payAmount
		} else if user.Wallet > 0 && user.Wallet < payAmount {
			walletPayAmount = user.Wallet
			wxPayAmount = payAmount - walletPayAmount
		}
	} else {
		wxPayAmount = payAmount
	}

	// 有余额付款则生成余额付款单
	if walletPayAmount > 0 {
		resp.WalletPaymentOrder, err = p.makePaymentOrder(tx, walletPayAmount, userId, orderId, model.PMOModeWallet)
		if err != nil {
			return response.PaymentOrderResp{}, err
		}
	}

	// 有微信支付则生成微信付款单
	if wxPayAmount > 0 {
		resp.WxPaymentOrder, err = p.makePaymentOrder(tx, wxPayAmount, userId, orderId, model.PMOModeWechat)
		if err != nil {
			return response.PaymentOrderResp{}, err
		}

		// 生成微信支付的参数
		resp.WxPayResp, err = wechat.NewPayment().GetPrepayBill(user.OpenID, "description", resp.WxPaymentOrder.PaymentOrderNo, wxPayAmount)
		if err != nil {
			return response.PaymentOrderResp{}, err
		}
	}

	// 没有微信支付的话，那就直接把订单改成已支付
	if wxPayAmount == 0 && walletPayAmount >= 0 {
		err = OrderService.PaySuccess(tx, orderId)
		if err != nil {
			return response.PaymentOrderResp{}, err
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
func (p *paymentService) RefundWalletPay(paymentOrder model.PaymentOrder, userId, amount int32) (
	refundOrder model.RefundOrder, err error) {
	// 获取用户信息
	user := model.User{}
	p.db.Where("user_id = ?", userId).First(&user)

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
	for _, v := range order.PaymentOrderList {
		// 付款金额如果大于退款金额，则退掉就结束
		// 付款金额如果
		payAmount := v.Amount

		if amount <= 0 {
			break
		}

		refundAmount := amount

		// 如果支付金额小于退款金额，则这一单整单退，然后剩余的钱下一单给退
		if payAmount < refundAmount {
			refundAmount = payAmount
			amount = amount - refundAmount
		}

		// 退款
		if v.PayMode == model.PMOModeWallet {
			// 退还到钱包
			_, err := p.RefundWalletPay(v, order.UserID, refundAmount)
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

	return
}

func (p *paymentService) PayNotify(c *gin.Context) (err error) {
	res, request, err := wechat.NewPayment().GetPayResult(c)
	if err != nil {
		return
	}

	paymentOrder := model.PaymentOrder{}
	if err = p.db.Where("payment_order_no = ?", res.OutTradeNo).First(&paymentOrder).Error; err != nil {
		log.GetLogger().Error("pay_error", zap.String("msg", "查询支付订单失败："+err.Error()))
		return
	}

	paymentOrder.NotifyID = request.ID
	paymentOrder.Resource = ""
	paymentOrder.BankType = *res.BankType
	paymentOrder.TransactionID = *res.TransactionId
	paymentOrder.TradeState = *res.TradeState
	if *res.TradeState == "SUCCESS" {
		paymentOrder.Status = model.PMOStatusSuccess
	}

	if err = p.db.Save(&paymentOrder).Error; err != nil {
		log.GetLogger().Error("pay_error", zap.String("msg", "更新支付订单失败："+err.Error()))
		return
	}

	if paymentOrder.OrderType == model.POTypeTable {
		// 开台订单回调
		//_, err = TableOrderService.PaySuccess(paymentOrder.OrderID, p.db)
	} else if paymentOrder.OrderType == model.POTypeRecharge {
		// 充值订单回调
		//_, err = RechargeOrderService.PaySuccess(paymentOrder.OrderID)
	}

	if err != nil {
		log.GetLogger().Error("pay_error", zap.String("msg", "更新订单失败："+err.Error()))
		return err
	}

	return
}

func (p *paymentService) DoOrderSuccess(db *gorm.DB, paymentOrder *model.PaymentOrder) (err error) {
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

	//err = p.DoWechatOrderSuccess(tx, &paymentOrder)
	//if err != nil {
	//	return err
	//}

	// 调用订单支付成功的方法
	err = OrderService.PaySuccess(tx, paymentOrder.OrderID)
	if err != nil {
		log.GetLogger().Error("pay_error", zap.String("msg", "更新订单失败："+err.Error()))
		return err
	}

	tx.Commit()

	return
}

// 将钱包余额支付的付款单改为支付成功
func (p *paymentService) DoWalletOrderSuccess(db *gorm.DB, paymentOrder *model.PaymentOrder) (err error) {
	user := model.User{}
	db.Where("user_id = ?", paymentOrder.UserID).First(&user)

	if user.Wallet < paymentOrder.Amount {
		log.GetLogger().Error("wallet_error", zap.String("msg", "余额小于支付金额，支付失败"))
		err = errors.New("余额小于支付金额，支付失败")
		return
	}

	// 更新付款单状态
	paymentOrder.Status = model.PMOStatusSuccess
	db.Save(&paymentOrder)

	// 扣除用户余额
	user.Wallet = user.Wallet - paymentOrder.Amount
	db.Save(&user)

	return
}

// 创建优惠券订单
func (p *paymentService) MakeCouponOrder(db *gorm.DB, amount int32, orderType int, orderId, userId int32, description string) (
	payment *jsapi.PrepayWithRequestPaymentResponse, err error) {

	user, _ := UserService.GetByUserId(userId)

	paymentOrder := model.PaymentOrder{
		PaymentOrderNo: tool.GenerateOrderNum(),
		OrderID:        orderId,
		Amount:         amount,
		OrderType:      orderType,

		PayMode: model.PMOModeWechat,
		Status:  model.PMOStatusDefault,
	}

	if err = db.Create(&paymentOrder).Error; err != nil {
		return
	}

	payment, err = wechat.NewPayment().GetPrepayBill(user.OpenID, description, paymentOrder.PaymentOrderNo, amount)
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
	}

	if err = db.Create(&order).Error; err != nil {
		return
	}

	return
}

// 创建预支付订单并生成预支付参数
func (p *paymentService) MakeWechatPrepayOrder(userId, amount int32, orderType int, orderId int32, description string) (
	payment *jsapi.PrepayWithRequestPaymentResponse, err error) {

	user, _ := UserService.GetByUserId(userId)

	// 写入payment_order表
	paymentOrder := model.PaymentOrder{
		PaymentOrderNo: tool.GenerateOrderNum(),
		OrderID:        orderId,
		Amount:         amount,
		OrderType:      orderType,
		PayMode:        model.PMOModeWechat,
		Status:         model.PMOStatusDefault,
	}

	if err = p.db.Create(&paymentOrder).Error; err != nil {
		return
	}

	payment, err = wechat.NewPayment().GetPrepayBill(user.OpenID, description, paymentOrder.PaymentOrderNo, amount)
	if err != nil {
		return
	}

	return
}
