package api

import (
	"fmt"
	"gin-api/request"
	"gin-api/response"
	"gin-api/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type taskApi struct{}

var TaskApi = new(taskApi)

func (*taskApi) List(c *gin.Context) {
	list := service.TaskService.List()

	response.Success(c, list)
}

func (*taskApi) ChangeStatus(c *gin.Context) {
	form := request.ChangeTaskStatus{}

	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	task, err := service.TaskService.ChangeStatus(form)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, task)
}

func (*taskApi) Save(c *gin.Context) {
	form := request.SaveTask{}
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	task, err := service.TaskService.Save(form)

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, task)
}

func (*taskApi) Detail(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Query("task_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
	}

	task, err := service.TaskService.Detail(int32(taskId))

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, task)
}

func (*taskApi) Log(c *gin.Context) {
	form := request.TaskLogList{}
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	fmt.Println(form)

	resp, err := service.TaskService.Log(form)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, resp)
}

func (*taskApi) Execute(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Query("task_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
	}

	err = service.TaskService.Execute(int32(taskId))

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, nil)
}

func (*taskApi) Stop(c *gin.Context) {
	logId, err := strconv.Atoi(c.Query("task_log_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
	}

	err = service.TaskService.StopTask(int64(logId))

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, nil)
}
