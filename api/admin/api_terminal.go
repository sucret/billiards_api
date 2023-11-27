package api

import (
	"billiards/request"
	"billiards/response"
	"billiards/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

type terminalApi struct{}

var TerminalApi = new(terminalApi)

func (*terminalApi) ChangeStatus(c *gin.Context) {
	var form request.ChangeTerminalStatus
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	terminal, err := service.TerminalService.ChangeStatus(form)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, terminal)
}

func (*terminalApi) Save(c *gin.Context) {
	var form request.SaveTerminal
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	terminal, err := service.TerminalService.Save(form)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, terminal)
}

func (*terminalApi) Delete(c *gin.Context) {
	terminalId, err := strconv.Atoi(c.Query("terminal_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}

	err = service.TerminalService.Delete(terminalId)
	if err != nil {
		return
	}
}
