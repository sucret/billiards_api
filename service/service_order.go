package service

import (
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	"billiards/pkg/tool"
	"billiards/request"
	"billiards/response"
	"errors"
	"gorm.io/gorm"
)

type orderService struct {
	db *gorm.DB
}

var OrderService = &orderService{
	db: mysql.GetDB(),
}

// 订单详情
type OrderInfo struct {
	model.Order
	TableOrder       model.TableOrder
	CouponOrder      model.CouponOrder
	RechargeOrder    model.RechargeOrder
	PaymentOrderList []model.PaymentOrder
}

// 统一获取订单支付状态
func (o *orderService) PayResult(orderId, userId int32) (resp response.PayResult, err error) {
	order := model.Order{}
	if err = o.db.Where("order_id = ? AND user_id = ?", orderId, userId).
		First(&order).Error; err != nil {
		err = errors.New("订单不存在")
		return
	}

	resp.Succeed = order.Status == model.OrderStatusPaySuccess

	user, _ := UserService.GetByUserId(userId)
	resp.Wallet = user.Wallet

	return
}

// 统一下单入口
func (o *orderService) Create(userId int32, param request.OrderCreate) (resp response.OrderResp, err error) {
	resp.Order.UserID = userId

	tx := o.db.Begin()
	defer tx.Rollback()

	// 1、创建主订单
	if err = tx.Save(&resp.Order).Error; err != nil {
		err = errors.New("创建订单失败，请重试")
		return
	}

	var totalPayAmount int32 = 0

	// 2、创建子订单
	//  创建球桌订单
	if param.TableID > 0 {
		tableOrder, err := TableOrderService.Create(tx, param.TableID, userId, resp.Order.OrderID, param.CouponID, param.UserCouponID)
		if err != nil {
			return response.OrderResp{}, err
		}

		tool.Dump(tableOrder)
		totalPayAmount = tableOrder.PayAmount + totalPayAmount
	}

	//  创建优惠券订单
	if param.CouponID > 0 {
		couponOrder, err := CouponOrderService.Create(tx, param.CouponID, userId, resp.Order.OrderID)
		if err != nil {
			return response.OrderResp{}, err
		}

		tool.Dump(couponOrder)
		totalPayAmount = couponOrder.PayAmount + totalPayAmount
	}

	//  创建充值订单
	if param.IsRecharge {
		rechargeOrder, err := RechargeOrderService.Create(tx, param.RechargeAmountType, userId, resp.Order.OrderID)
		if err != nil {
			return response.OrderResp{}, err
		}

		totalPayAmount = rechargeOrder.PayAmount + totalPayAmount
	}

	// 3、计算订单金额，生成支付参数
	paymentResp, err := PaymentService.CreateOrder(tx, userId, resp.Order.OrderID, totalPayAmount, !param.IsRecharge)
	if err != nil {
		return response.OrderResp{}, err
	}

	resp.WxPayResp = paymentResp.WxPayResp

	// 标记是否需要微信支付
	if paymentResp.WxPaymentOrder.OrderID > 0 {
		resp.NeedWxPay = true
	}

	tx.Commit()

	return
}

// 支付成功
func (o *orderService) PaySuccess(db *gorm.DB, orderId int32) (err error) {
	order, err := o.getOrderInfo(db, orderId)
	if order.OrderID == 0 || order.Status != 1 {
		err = errors.New("订单信息不存在")
		return
	}

	// 子订单支付成功
	// 充值订单支付成功
	if order.RechargeOrder.RechargeOrderID > 0 {
		err = RechargeOrderService.PaySuccess(db, &order.RechargeOrder)
		if err != nil {
			return
		}
	}

	// 优惠券订单支付成功
	if order.CouponOrder.CouponOrderID > 0 {
		err = CouponOrderService.PaySuccess(db, &order.CouponOrder)
		if err != nil {
			return
		}
	}

	// 球桌/棋牌桌 订单支付成功
	// 需要放在优惠券之后执行，因为这一单使用了组合支付，优惠券支付先成功之后才会用在这一单上边
	if order.TableOrder.TableOrderID > 0 {
		err = TableOrderService.PaySuccess(db, &order.TableOrder)
		if err != nil {
			return
		}
	}

	// 主订单支付成功
	order.Status = 2
	err = db.Save(&order.Order).Error

	for _, v := range order.PaymentOrderList {
		err = PaymentService.DoOrderSuccess(db, &v)
		if err != nil {
			return
		}
	}

	return
}

func (o *orderService) getOrderInfo(db *gorm.DB, orderId int32) (order OrderInfo, err error) {
	if err = db.Where("order_id = ?", orderId).First(&order.Order).Error; err != nil {
		err = errors.New("订单不存在")
		return
	}

	db.Where("order_id = ?", orderId).First(&order.RechargeOrder)
	db.Where("order_id = ?", orderId).First(&order.CouponOrder)
	db.Where("order_id = ?", orderId).First(&order.TableOrder)
	db.Where("order_id = ?", orderId).First(&order.PaymentOrderList)

	return
}
