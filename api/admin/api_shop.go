package api

import (
	"billiards/request"
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type shopApi struct{}

var ShopApi = new(shopApi)

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

func (*shopApi) List(c *gin.Context) {
	list := service.ShopService.List()
	response.Success(c, list)
}

func (*shopApi) Save(c *gin.Context) {
	var form request.SaveShop
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}
	shop, err := service.ShopService.Save(form)

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, shop)
}
