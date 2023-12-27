package router

import (
	api "billiards/api/admin"
	"billiards/middleware"
	"billiards/service"

	"github.com/gin-gonic/gin"
)

func setAdminRouter(r *gin.Engine) {

	// 登陆接口不需要走jwt验证
	r.POST("/admin/login", api.AdminApi.Login)
	r.POST("/admin/get-login-sms", api.AdminApi.AdminSendLoginSms)

	// 菜单接口不走权限验证
	r.GET("/admin/admin/menu",
		middleware.JWTAuth(service.AppGuardName),
		gin.Logger(),
		middleware.CustomRecovery(),
		api.AdminApi.MenuList)

	adminRouter := r.Group("/admin").Use(
		middleware.JWTAuth(service.AppGuardName),
		middleware.CheckAdminPermission(),
		gin.Logger(),
		middleware.CustomRecovery())
	{
		adminRouter.GET("/profile", api.AdminApi.Profile)
		adminRouter.GET("/admin/list", api.AdminApi.List)
		adminRouter.GET("/admin/detail", api.AdminApi.Detail)
		adminRouter.POST("/admin/save", api.AdminApi.Save)

		// 角色
		adminRouter.GET("/role/list", api.RoleApi.List)
		adminRouter.POST("/role/save", api.RoleApi.Save)
		adminRouter.GET("/role/detail", api.RoleApi.Detail)
		adminRouter.GET("/role/node", api.RoleApi.RoleNode)
		adminRouter.POST("/role/node/save", api.RoleApi.SaveRoleNode)

		// 权限
		adminRouter.GET("/node/list", api.NodeApi.List)
		adminRouter.POST("/node/save", api.NodeApi.Save)
		adminRouter.GET("/node/tree", api.NodeApi.NodeTree)
		adminRouter.GET("/node/detail", api.NodeApi.Detail)
		adminRouter.GET("/node/delete", api.NodeApi.Delete)

		// 定时任务
		adminRouter.GET("/task/list", api.TaskApi.List)
		adminRouter.POST("/task/change-status", api.TaskApi.ChangeStatus)
		adminRouter.POST("/task/save", api.TaskApi.Save)
		adminRouter.GET("/task/detail", api.TaskApi.Detail)
		adminRouter.POST("/task/log", api.TaskApi.Log)
		adminRouter.GET("/task/execute", api.TaskApi.Execute)
		adminRouter.GET("/task/stop", api.TaskApi.Stop)

		// 数据库操作
		adminRouter.GET("/mysql/table-list", api.MysqlApi.Tables)
		adminRouter.POST("/mysql/execute-sql", api.MysqlApi.Execute)

		// Redis缓存操作
		adminRouter.POST("/cache/redis-query", api.RedisApi.Query)

		// 店铺
		adminRouter.GET("/shop/list", api.ShopApi.List)
		adminRouter.POST("/shop/save", api.ShopApi.Save)
		adminRouter.GET("/shop/detail", api.ShopApi.Detail)

		// 终端
		adminRouter.POST("/terminal/change-status", api.TerminalApi.ChangeStatus)
		adminRouter.POST("/terminal/save", api.TerminalApi.Save)
		adminRouter.GET("/terminal/delete", api.TerminalApi.Delete)

		adminRouter.GET("/table/enable", api.TableApi.Enable)
		adminRouter.POST("table/save", api.TableApi.Save)

		// 订单
		adminRouter.POST("/order/list", api.OrderApi.List)
		adminRouter.GET("/order/detail", api.OrderApi.Detail)

		// 用户
		adminRouter.POST("/user/list", api.UserApi.List)
	}
}
