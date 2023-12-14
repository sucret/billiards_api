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

func (o *orderService) List(form request.OrderList) (resp response.OrderListResp, err error) {
	query := o.db.Preload("TableOrder").
		Preload("TableOrder.Table").
		Preload("TableOrder.Table.Shop").
		Preload("PaymentOrderList").
		Preload("CouponOrder").
		Preload("CouponOrder.Coupon")

	query.Model(&[]model.Order{}).Count(&resp.Total)

	err = query.Order("order_id DESC").Offset((form.Page - 1) * form.PageSize).
		Limit(form.PageSize).Find(&resp.List).Error

	return
}

// 续费
// 由于当前业务单据只有开台订单可以续费，所以这里省掉了按业务单类型续费逻辑
// 后续如果有其他类型的业务单据续费则需要按类别处理
func (o *orderService) Renewal(orderId, userId int32) (resp response.OrderResp, err error) {
	tx := o.db.Begin()
	defer tx.Rollback()

	// 获取订单信息
	order, _ := o.GetOrderInfo(tx, orderId)
	if order.UserID != userId {
		err = errors.New("订单不存在")
		return
	}

	resp.Order = order.Order

	// 获取需要支付的金额
	payAmount, err := TableOrderService.GetRenewalAmount(order.TableOrder)
	if err != nil {
		return
	}

	// 生成支付参数
	paymentResp, err := PaymentService.createOrder(tx, userId, order.OrderID, payAmount, true, true)
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

// 订单详情
func (o *orderService) Detail(orderId, userId int32) (detail response.OrderDetailResp, err error) {
	order, _ := o.GetOrderInfo(o.db, orderId)
	if userId != -1 && order.UserID != userId {
		err = errors.New("订单不存在")
		return
	}

	tool.Dump(order)
	detail.Order = order.Order
	detail.TableOrder.TableOrder = order.TableOrder
	detail.PaymentOrderList = order.PaymentOrderList
	detail.CouponOrder = order.CouponOrder

	TableOrderService.formatTableOrder(&detail)

	return
}

// 终止订单
func (o *orderService) Terminate(orderId, userId int32) (err error) {
	// 1、获取订单信息
	tx := o.db.Begin()
	defer tx.Rollback()

	order, err := o.GetOrderInfo(tx, orderId)

	if order.OrderID == 0 || order.Status != model.OrderStatusPaySuccess {
		err = errors.New("订单信息不存在")
		return
	}

	if order.UserID != userId {
		err = errors.New("用户只能终止自己的订单")
		return
	}

	var orderAmount int32 = 0
	var couponAmount int32 = 0
	var tableOrderAmount int32 = 0

	// 2、计算没个类型的订单现在应该付多少钱
	// 计算开台订单的退款金额
	if order.TableOrder.OrderID > 0 {
		tableOrderAmount, err = TableOrderService.settlement(&order.TableOrder)
		if err != nil {
			return
		}
		orderAmount = orderAmount + tableOrderAmount
	}

	// 结算优惠券金额，如果优惠券未使用，则退掉优惠券
	if order.CouponOrder.OrderID > 0 {
		couponAmount, err = CouponOrderService.settlement(&order.CouponOrder)
		if err != nil {
			return
		}
		orderAmount = orderAmount + couponAmount
	}

	// 3、终止各个业务订单
	if err = TableOrderService.terminate(tx, &order.TableOrder); err != nil {
		return
	}

	if couponAmount > 0 {
		if err = CouponOrderService.terminate(tx, &order.CouponOrder); err != nil {
			return
		}
	}

	// 4、修改主订单的状态
	order.Status = model.OrderStatusFinised
	err = tx.Save(&order.Order).Error
	if err != nil {
		return
	}

	tx.Commit()

	// 5、计算退款金额，如果金额大于0则操作退款
	refundAmount := order.TableOrder.PayAmount + order.CouponOrder.PayAmount - orderAmount
	if refundAmount > 0 {
		PaymentService.refund(refundAmount, &order.PaymentOrderList)
	}

	return
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
		tableOrder, err := TableOrderService.create(tx, param.TableID, userId, resp.Order.OrderID, param.CouponID, param.UserCouponID)
		if err != nil {
			return response.OrderResp{}, err
		}

		totalPayAmount = tableOrder.PayAmount + totalPayAmount
	}

	//  创建优惠券订单
	if param.CouponID > 0 {
		couponOrder, err := CouponOrderService.create(tx, param.CouponID, userId, resp.Order.OrderID)
		if err != nil {
			return response.OrderResp{}, err
		}

		totalPayAmount = couponOrder.PayAmount + totalPayAmount
	}

	//  创建充值订单
	if param.IsRecharge {
		rechargeOrder, err := RechargeOrderService.create(tx, param.RechargeAmountType, userId, resp.Order.OrderID)
		if err != nil {
			return response.OrderResp{}, err
		}

		totalPayAmount = rechargeOrder.PayAmount + totalPayAmount
	}

	// 3、计算订单金额，生成支付参数
	paymentResp, err := PaymentService.createOrder(tx, userId, resp.Order.OrderID, totalPayAmount, !param.IsRecharge, false)
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
	order, err := o.GetOrderInfo(db, orderId)
	if order.OrderID == 0 || order.Status != 1 {
		err = errors.New("订单信息不存在")
		return
	}

	// 子订单支付成功
	// 充值订单支付成功
	if order.RechargeOrder.RechargeOrderID > 0 {
		err = RechargeOrderService.paySuccess(db, &order.RechargeOrder)
		if err != nil {
			return
		}
	}

	// 优惠券订单支付成功
	if order.CouponOrder.CouponOrderID > 0 {
		err = CouponOrderService.paySuccess(db, &order.CouponOrder)
		if err != nil {
			return
		}
	}

	// 球桌/棋牌桌 订单支付成功
	// 需要放在优惠券之后执行，因为这一单使用了组合支付，优惠券支付先成功之后才会用在这一单上边
	if order.TableOrder.TableOrderID > 0 {
		err = TableOrderService.paySuccess(db, &order.TableOrder)
		if err != nil {
			return
		}
	}

	// 主订单支付成功
	order.Status = 2
	err = db.Save(&order.Order).Error

	// 付款单支付成功
	for _, v := range order.PaymentOrderList {
		err = PaymentService.doOrderSuccess(db, &v)
		if err != nil {
			return
		}
	}

	return
}

func (o *orderService) GetOrderInfo(db *gorm.DB, orderId int32) (order OrderInfo, err error) {
	if err = db.Where("order_id = ?", orderId).First(&order.Order).Error; err != nil {
		err = errors.New("订单不存在")
		return
	}

	db.Where("order_id = ?", orderId).First(&order.RechargeOrder)
	db.Where("order_id = ?", orderId).Preload("Coupon").First(&order.CouponOrder)
	db.Where("order_id = ?", orderId).Preload("Table").Preload("Table.Shop").First(&order.TableOrder)
	db.Where("order_id = ?", orderId).Preload("RefundOrderList").Find(&order.PaymentOrderList)

	return
}
