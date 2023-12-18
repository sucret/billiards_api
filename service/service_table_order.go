package service

import (
	"billiards/pkg/log"
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	redis_ "billiards/pkg/redis"
	"billiards/pkg/tool"
	"billiards/response"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"sync"
	"time"
)

type tableOrderService struct {
	db    *gorm.DB
	redis *redis.Client
	lock  sync.Mutex
}

var TableOrderService = &tableOrderService{
	db:    mysql.GetDB(),
	redis: redis_.GetRedis(),
	lock:  sync.Mutex{},
}

func (o *tableOrderService) GetRenewalAmount(order model.TableOrder) (amount int32, err error) {
	if order.Table.Shop.ShopID > 0 {
		amount = order.Table.Shop.Deposit
	} else {
		err = errors.New("未获取到店铺信息")
	}

	return
}

func (o *tableOrderService) GetByTable(tableId, userId int) (order model.TableOrder, err error) {
	if err = o.db.Where("table_id = ? AND user_id = ? AND status = ?",
		tableId, userId, model.OrderStatusPaySuccess).
		First(&order).Error; err != nil {
		return
	}
	return
}

func (o *tableOrderService) formatTableOrder(detail *response.OrderDetailResp) {
	var totalAmount, orderAmount int32
	for _, v := range detail.PaymentOrderList {
		totalAmount = totalAmount + v.Amount
	}
	orderAmount = totalAmount

	tool.Dump(detail.PaymentOrderList)
	fmt.Println(totalAmount)

	// 如果有优惠券，则需要加上优惠券的时长
	var couponDuration int32 = 0

	if detail.CouponOrder.CouponOrderID > 0 {
		coupon, err := CouponService.GetById(detail.CouponOrder.CouponID)
		tool.Dump(coupon)
		if err != nil {
			return
		}

		orderAmount = orderAmount - coupon.Price
		couponDuration = coupon.Duration
	}

	// 使用时长
	if detail.TableOrder.Status == model.OrderStatusPaySuccess {
		detail.TableOrder.UsedMinutes = int32(time.Now().Sub(time.Time(detail.TableOrder.StartedAt)).Minutes())
	} else if detail.TableOrder.Status == model.OrderStatusFinised {
		detail.TableOrder.UsedMinutes =
			int32(time.Time(detail.TableOrder.TerminatedAt).Sub(time.Time(detail.TableOrder.StartedAt)).Minutes())
	}

	// 总时长 = 优惠券时长 + 订单时长
	// 订单时长 = 订单金额 / 单价 * 60
	// 订单金额 = 总金额 - 优惠券金额
	detail.TableOrder.TotalMinutes = int32(math.Ceil(float64(orderAmount)/float64(detail.TableOrder.Table.Price)*60)) + couponDuration

	// 剩余时长
	detail.TableOrder.RemainMinutes = detail.TableOrder.TotalMinutes - detail.TableOrder.UsedMinutes
	if detail.TableOrder.RemainMinutes < 0 {
		detail.TableOrder.RemainMinutes = 0
	}

	// 待结算金额
	// 使用时长 - 优惠券时长
	waitSettlementMinutes := math.Max(float64(detail.TableOrder.UsedMinutes-couponDuration), 0)
	detail.TableOrder.WaitSettlementAmount = int32(float64(detail.TableOrder.Table.Price) / 60 * waitSettlementMinutes)
}

// 终止开台订单
func (o *tableOrderService) terminate(db *gorm.DB, tableOrder *model.TableOrder) (err error) {
	// 关闭球桌
	_, err = TableService.Disable(db, tableOrder.TableID)
	if err != nil {
		return
	}

	// 修改开台订单状态
	tableOrder.Status = model.OrderStatusFinised
	if err = db.Save(&tableOrder).Error; err != nil {
		err = errors.New("关闭失败")
		return
	}

	return
}

// 结算订单，查看订单现在有多少钱
// 返回的是减掉优惠券抵扣的时长
func (o *tableOrderService) settlement(tableOrder *model.TableOrder) (amount int32, err error) {
	// 结束时间
	tableOrder.TerminatedAt = model.Time(time.Now())

	// 时长(分钟）
	minutes := int32(time.Time(tableOrder.TerminatedAt).Sub(time.Time(tableOrder.StartedAt)).Minutes())

	// 查看是否有用优惠券，如果用了，那就需要把优惠券的时长去掉再计算金额
	if tableOrder.CouponID > 0 {
		coupon, err := CouponService.GetById(tableOrder.CouponID)
		if err != nil {
			return 0, err
		}

		minutes = minutes - coupon.Duration
		if minutes < 0 {
			minutes = 0
		}
	}

	// 金额
	tableOrder.Amount = int32(float64(tableOrder.Table.Price) / float64(60) * float64(minutes))

	amount = tableOrder.Amount

	return
}

// 终止订单
//
// 检查球杆是否归还
// 检查球是否够
// 关闭设备
// 结算
// 发起微信退款
// 修改订单状态
func (o *tableOrderService) Terminate(userId, tableOrderId int) (order *response.OrderDetail, err error) {
	tx := o.db.Begin()

	order = &response.OrderDetail{}
	// 获取订单信息
	if err = tx.Set("gorm:query_option", "FOR UPDATE").
		Preload("Table").
		Where("user_id = ? AND table_order_id = ? AND status = ?",
			userId, tableOrderId, model.OrderStatusPaySuccess).
		First(&order).Error; err != nil {
		err = errors.New("未查询到订单信息")
		tx.Rollback()
		return
	}

	orderDetail, err := OrderService.GetOrderInfo(tx, order.OrderID)
	if err != nil {
		return
	}

	// 关闭球桌
	// todo
	_, err = TableService.Disable(tx, orderDetail.TableOrder.TableID)
	if err != nil {
		tx.Rollback()
		return
	}

	//tool.Dump(order)
	//tx.Rollback()
	//return

	// 修改订单状态
	order.Status = model.OrderStatusFinised
	order.TerminatedAt = model.Time(time.Now())

	// 结算 回写订单金额
	//o.settlement(&orderDetail.TableOrder)

	// todo 退押金
	o.refund(orderDetail.TableOrder)

	if err = tx.Save(&order.TableOrder).Error; err != nil {
		tx.Rollback()
		err = errors.New("关闭失败")
		return
	}

	tx.Commit()

	// 重新获取订单信息
	order, err = o.Detail(userId, tableOrderId)

	return
}

// 退款
func (o *tableOrderService) refund(order model.TableOrder) {
	var totalAmount int32
	//for _, v := range order.PaymentOrderList {
	//	totalAmount = totalAmount + v.Amount
	//}

	refundAmount := totalAmount - order.Amount

	PaymentService.RefundOrder(order, refundAmount)
}

func (o *tableOrderService) Detail(userId, orderId int) (order *response.OrderDetail, err error) {
	if err = o.db.Preload("Table").
		Preload("Table.Shop").
		//Preload("PaymentOrderList", func(db *gorm.DB) *gorm.DB {
		//	return db.Select("order_id,amount").
		//		Where("status = ?", model.PMOStatusSuccess).
		//		Order("pay_mode")
		//}).
		Where("user_id = ? AND order_id = ?", userId, orderId).
		First(&order).Error; err != nil {
		err = errors.New("未查询到订单信息")
	}

	o.formatClientOrder(order)

	return
}

// 计算订单的金额和时间
func (o *tableOrderService) formatClientOrder(order *response.OrderDetail) {
	var totalAmount int32
	//for _, v := range order.PaymentOrderList {
	//	totalAmount = totalAmount + v.Amount
	//}

	// 总时长 （当前支付的金额 / 单价 * 60）
	order.TotalMinutes = (totalAmount / order.Table.Price) * 60

	if order.Status == model.OrderStatusPaySuccess {
		// 使用时长
		order.UsedMinutes = int32(time.Now().Sub(time.Time(order.StartedAt)).Minutes())
	} else if order.Status == model.OrderStatusFinised {
		order.UsedMinutes = int32(time.Time(order.TerminatedAt).Sub(time.Time(order.StartedAt)).Minutes())
	}

	// 剩余时长
	order.RemainMinutes = order.TotalMinutes - order.UsedMinutes
	if order.RemainMinutes < 0 {
		order.RemainMinutes = 0
	}

	//tool.Dump(order)
}

func (o *tableOrderService) List(userId int, orderType int) (list []*response.TableOrderDetail, err error) {
	query := o.db.Preload("Table").Preload("Table.Shop")

	if orderType == 0 {
		query.Where("status in ? AND user_id = ?",
			[]int{model.OrderStatusPaySuccess, model.OrderStatusFinised}, userId)
	} else {
		query.Where("status = ? AND user_id = ?", orderType, userId)
	}
	if err = query.Order("order_id desc").
		Find(&list).Error; err != nil {
	}

	for _, v := range list {
		if time.Time(v.TerminatedAt).IsZero() {
			v.Duration = int32(time.Now().Sub(time.Time(v.StartedAt)).Minutes())
		} else {
			v.Duration = int32(time.Time(v.TerminatedAt).Sub(time.Time(v.StartedAt)).Minutes())
		}
	}

	return
}

// 查询用户自己订单的支付结果
func (o *tableOrderService) GetPayStatus(orderId int, userId int32) (successful bool, err error) {
	order := model.TableOrder{}
	if err = o.db.Where("order_id = ? AND user_id = ?", orderId, userId).
		First(&order).Error; err != nil {
		err = errors.New("订单不存在")
		return
	}

	successful = order.Status == model.OrderStatusPaySuccess

	return
}

// 取消10分钟前未支付的订单
func (o *tableOrderService) TimingCancel() {
	var orderList []model.TableOrder

	t, _ := time.ParseDuration("-10m")
	o.db.Where("status = ? AND created_at < ?", model.OrderStatusDefault,
		time.Now().Add(t).Format("2006-01-02 15:04:05")).
		Find(&orderList)

	for _, v := range orderList {
		v.Status = model.OrderStatusAutoCancel
		err := o.db.Save(&v).Error
		if err != nil {
			log.GetLogger().Error("auto cancel error",
				zap.Any("order", v),
				zap.String("error", err.Error()))
		}
	}
}

// 创建开台子订单
func (o *tableOrderService) create(tx *gorm.DB, tableId, userId, orderId, couponId, userCouponId int32) (
	tableOrder model.TableOrder, err error) {

	// todo 暂时先在创建订单的时候关闭过期的订单，后边需要移动到定时任务里边去
	go o.TimingCancel()

	// 1、判断球桌是否可用
	table := model.Table{}
	if err = tx.Preload("Shop").
		Where("table_id = ? AND status = ?", tableId, model.TableStatusClose).
		First(&table).Error; err != nil {

		err = errors.New("当前球桌不可用，请更换其它球桌")
		return
	}

	// 2、判断当前球桌是否有其他人正在支付中的订单
	tmpOrder := model.TableOrder{}
	tx.Where("table_id = ? AND status = ? AND user_id !=", tableId, model.OrderStatusDefault, userId).
		First(&tmpOrder)
	if tmpOrder.OrderID > 0 {
		err = errors.New("当前球桌不可用，请更换其它球桌")
		return
	}

	// 3、创建开台子订单
	tableOrder = model.TableOrder{
		OrderID:      orderId,
		ShopID:       table.ShopID,
		TableID:      table.TableID,
		Status:       model.OrderStatusDefault,
		UserID:       userId,
		CouponID:     couponId,
		UserCouponID: userCouponId,
		PayAmount:    table.Shop.Deposit,
	}
	if err = tx.Create(&tableOrder).Error; err != nil {
		log.GetLogger().Error("create_order", zap.String("err_msg", err.Error()),
			zap.Any("order", tableOrder))
		err = errors.New("创建订单，请重试")
	}

	return
}

func (o *tableOrderService) paySuccess(db *gorm.DB, order *model.TableOrder) (err error) {

	order.Status = model.OrderStatusPaySuccess
	order.StartedAt = model.Time(time.Now())

	// 优惠券处理
	if order.CouponID > 0 {
		// 有优惠券但是没绑定，说明是这一单购买的
		userCoupon, err := UserCouponService.UseByCouponID(db, order.CouponID, order.UserID)
		if err != nil {
			return err
		}
		order.UserCouponID = userCoupon.UserCouponID
	} else if order.UserCouponID > 0 {
		// 有绑定优惠券，说明是提前购买的
		userCoupon, err := UserCouponService.Use(db, order.UserCouponID, order.UserID)
		if err != nil {
			return err
		}
		order.CouponID = userCoupon.CouponID
	}

	if err = db.Save(&order).Error; err != nil {
		log.GetLogger().Error("pay_notify",
			zap.String("msg", "更新订单状态失败"),
			zap.String("err_msg", err.Error()),
			zap.Int32("order_id", order.TableOrderID))
	}

	// 开台
	table, err := TableService.Enable(db, order.TableID)
	if err != nil {
		log.GetLogger().Error("pay_notify",
			zap.String("msg", "开台失败"),
			zap.Any("table", table),
			zap.Int32("order_id", order.TableOrderID))

		return err
	}

	return
}
