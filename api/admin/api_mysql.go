package api

import (
	"gin-api/request"
	"gin-api/response"
	"gin-api/service"
	"github.com/gin-gonic/gin"
)

type mysqlApi struct{}

var MysqlApi = new(mysqlApi)

// 查看数据列表
func (*mysqlApi) Tables(c *gin.Context) {
	list, _ := service.MysqlService.Tables()

	response.Success(c, list)
}

// 执行sql
func (*mysqlApi) Execute(c *gin.Context) {
	var form request.ExecuteSql
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	data, err := service.MysqlService.Execute(form.Sql)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, data)
}
