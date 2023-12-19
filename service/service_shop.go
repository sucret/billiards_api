package service

import (
	"billiards/pkg/mysql"
	"billiards/pkg/mysql/model"
	redis_ "billiards/pkg/redis"
	"billiards/pkg/tool"
	"billiards/request"
	"billiards/response"
	"errors"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	"math"
)

type shopService struct {
	db    *gorm.DB
	redis *redis.Client
}

var ShopService = &shopService{
	db:    mysql.GetDB(),
	redis: redis_.GetRedis(),
}

func (s *shopService) Status(shopId int) (resp response.ShopStatusResp, err error) {
	shop := model.Shop{}
	err = s.db.Preload("TableList").
		Where("shop_id = ?", shopId).
		First(&shop).Error

	for _, v := range shop.TableList {
		resp.TableStatusList = append(resp.TableStatusList, response.TableStatus{TableID: v.TableID, Status: v.Status})
	}

	return
}

func (s *shopService) Detail(shopId int) (shop model.Shop, err error) {
	err = s.db.Preload("TableList").
		Preload("TableList.TerminalList").
		Where("shop_id = ?", shopId).
		First(&shop).Error

	if err != nil {
		err = errors.New("店铺不存在")
	}
	return
}

func (s *shopService) List() (list []*response.Shop) {
	s.db.Order("shop_id DESC").
		Preload("TableList").
		Find(&list)

	for _, v := range list {
		v.BilliardsTableNum = len(v.TableList)

		for _, t := range v.TableList {
			if t.Status == model.TableStatusClose {
				v.BilliardsTableFreeNum = v.BilliardsTableFreeNum + 1
			}
		}
		v.BilliardsPrice = 0
		if len(v.TableList) > 0 {
			v.BilliardsPrice = v.TableList[0].Price
		}
	}

	return
}

func (s *shopService) ListWithDistance(lat, lng float64) (shopList []*response.Shop) {
	shopList = s.List()

	if lat > 0 && lng > 0 {
		for _, v := range shopList {
			if v.Latitude > 0 && v.Longitude > 0 {
				v.Distance = math.Round(tool.Distance(lat, lng, v.Latitude, v.Longitude)/100) / 10
			}
		}
	}

	return
}

func (s *shopService) Save(form request.SaveShop) (shop model.Shop, err error) {

	if form.ShopID != 0 {
		if err = s.db.Where("shop_id = ?", form.ShopID).First(&shop).Error; err != nil {
			err = errors.New("店铺不存在")
			return
		}

		shop.Name = form.Name
		shop.Status = form.Status
		shop.Address = form.Address
		if err = s.db.Save(&shop).Error; err != nil {
			return
		}
	} else {
		shop = model.Shop{
			Name:    form.Name,
			Status:  form.Status,
			Address: form.Address,
		}
		err = s.db.Create(&shop).Error
	}
	return
}
