package api

import (
	"billiards/request"
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
)

type tableApi struct{}

var TableApi = new(tableApi)

func (*tableApi) Save(c *gin.Context) {
	var form request.SaveTable
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	table, err := service.TableService.Save(form)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, table)
}

func (*tableApi) Enable(c *gin.Context) {
	//tableId, err := strconv.Atoi(c.Query("table_id"))
	//if err != nil {
	//	response.BusinessFail(c, "参数错误")
	//	return
	//}
	//
	//table, err := service.TableService.Enable(int32(tableId))
	//if err != nil {
	//	response.BusinessFail(c, err.Error())
	//	return
	//}
	//
	//response.Success(c, table)
}
