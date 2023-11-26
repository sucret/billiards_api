package api

import (
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type shopApi struct{}

var ShopApi = new(shopApi)

func (*shopApi) List(c *gin.Context) {
	list := service.ShopService.List()
	response.Success(c, list)
}

func (*shopApi) Detail(c *gin.Context) {
	shopId, err := strconv.Atoi(c.Query("shop_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}

	shop, err := service.ShopService.Detail(shopId)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, shop)
}
