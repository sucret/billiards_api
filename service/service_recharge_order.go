package service

import (
	"billiards/pkg/config"
	"billiards/pkg/log"
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	redis_ "billiards/pkg/redis"
	"billiards/response"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

type rechargeOrderService struct {
	db    *gorm.DB
	redis *redis.Client
	lock  sync.Mutex
}

var RechargeOrderService = &rechargeOrderService{
	db:    mysql.GetDB(),
	redis: redis_.GetRedis(),
	lock:  sync.Mutex{},
}

func (r *rechargeOrderService) PayResult(orderId, userId int32) (successful bool, err error) {
	order := model.RechargeOrder{}
	if err = r.db.Where("order_id = ? AND user_id = ?", orderId, userId).
		First(&order).Error; err != nil {
		err = errors.New("订单不存在")
		return
	}

	successful = order.Status == model.RechargeOrderStatusPaid

	return
}

// 创建充值订单
func (r *rechargeOrderService) Create(userId int32, rechargeAmount int) (
	resp response.RechargeOrderPrePayParam, err error) {
	conf := config.GetConfig().RechargeAmount
	// 超出配置范围则按最低档处理
	if rechargeAmount < 0 || rechargeAmount >= len(conf) {
		rechargeAmount = 0
	}

	resp.Order = &model.RechargeOrder{
		UserID:        userId,
		Amount:        conf[rechargeAmount].Amount,
		BundledAmount: conf[rechargeAmount].BundledAmount,
		Status:        model.RechargeOrderStatusDefault,
	}

	// 创建订单
	if err = r.db.Create(resp.Order).Error; err != nil {
		log.GetLogger().Error("create_recharge_order_err", zap.String("msg", err.Error()))
		err = errors.New("创建订单失败，请重试")
		return
	}

	// 创建预付款单
	// 生成预支付的参数
	jsapiParam, err := PaymentService.MakePrepayOrder(
		userId, resp.Order.Amount, model.POTypeRecharge, resp.Order.OrderID, "会员充值")
	if err != nil {
		return
	}

	//resp.Order = &order
	resp.JsApi = jsapiParam

	return
}

// 订单支付成功回调
func (r *rechargeOrderService) PaySuccess(orderId int32) (order model.RechargeOrder, err error) {
	// 查询订单
	if err = r.db.Where("order_id = ?", orderId).Find(&order).Error; err != nil {
		fmt.Println(err)
		log.GetLogger().Error("recharge_notify_error", zap.String("msg", err.Error()))
		return
	}

	order.Status = model.RechargeOrderStatusPaid

	_, err = UserService.Recharge(order.Amount+order.BundledAmount, order.UserID)

	r.db.Save(&order)
	return
}
