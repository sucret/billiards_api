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

func (r *rechargeOrderService) PayResult(orderId, userId int32) (resp response.RechargeResult, err error) {
	order := model.RechargeOrder{}
	if err = r.db.Where("order_id = ? AND user_id = ?", orderId, userId).
		First(&order).Error; err != nil {
		err = errors.New("订单不存在")
		return
	}

	resp.Succeed = order.Status == model.RechargeOrderStatusPaid

	user, _ := UserService.GetByUserId(userId)
	resp.Wallet = user.Wallet

	return
}

func (r *rechargeOrderService) Create(tx *gorm.DB, rechargeAmountType int, userId, orderId int32) (
	rechargeOrder model.RechargeOrder, err error) {

	conf := config.GetConfig().RechargeAmount
	// 超出配置范围则按最低档处理
	if rechargeAmountType < 0 || rechargeAmountType >= len(conf) {
		rechargeAmountType = 0
	}
	rechargeOrder = model.RechargeOrder{
		OrderID:       orderId,
		UserID:        userId,
		PayAmount:     conf[rechargeAmountType].Amount,
		Amount:        conf[rechargeAmountType].Amount,
		BundledAmount: conf[rechargeAmountType].BundledAmount,
		Status:        model.RechargeOrderStatusDefault,
	}

	// 创建充值订单
	if err = tx.Create(&rechargeOrder).Error; err != nil {
		log.GetLogger().Error("create_recharge_order_err", zap.String("msg", err.Error()))
		err = errors.New("创建订单失败，请重试")
		return
	}
	fmt.Println("f2323")

	return
}

// 订单支付成功回调
func (r *rechargeOrderService) PaySuccess(db *gorm.DB, order *model.RechargeOrder) (err error) {
	// 查询订单
	//if err = r.db.Where("order_id = ?", orderId).Find(&order).Error; err != nil {
	//	fmt.Println(err)
	//	log.GetLogger().Error("recharge_notify_error", zap.String("msg", err.Error()))
	//	return
	//}

	// 修改订单状态
	order.Status = model.RechargeOrderStatusPaid
	db.Save(&order)

	// 给用户账户充值
	_, err = UserService.Recharge(db, order.Amount+order.BundledAmount, order.UserID)

	return
}
