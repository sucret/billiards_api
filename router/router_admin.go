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

	adminRouters := r.Group("/admin").Use(
		middleware.JWTAuth(service.AppGuardName),
		middleware.CheckAdminPermission(),
		gin.Logger(),
		middleware.CustomRecovery())
	{
		adminRouters.GET("/profile", api.AdminApi.Profile)
		adminRouters.GET("/admin/list", api.AdminApi.List)
		adminRouters.GET("/admin/detail", api.AdminApi.Detail)
		adminRouters.POST("/admin/save", api.AdminApi.Save)

		// 角色
		adminRouters.GET("/role/list", api.RoleApi.List)
		adminRouters.POST("/role/save", api.RoleApi.Save)
		adminRouters.GET("/role/detail", api.RoleApi.Detail)
		adminRouters.GET("/role/node", api.RoleApi.RoleNode)
		adminRouters.POST("/role/node/save", api.RoleApi.SaveRoleNode)

		// 权限
		adminRouters.GET("/node/list", api.NodeApi.List)
		adminRouters.POST("/node/save", api.NodeApi.Save)
		adminRouters.GET("/node/tree", api.NodeApi.NodeTree)
		adminRouters.GET("/node/detail", api.NodeApi.Detail)
		adminRouters.GET("/node/delete", api.NodeApi.Delete)

		// 定时任务
		adminRouters.GET("/task/list", api.TaskApi.List)
		adminRouters.POST("/task/change-status", api.TaskApi.ChangeStatus)
		adminRouters.POST("/task/save", api.TaskApi.Save)
		adminRouters.GET("/task/detail", api.TaskApi.Detail)
		adminRouters.POST("/task/log", api.TaskApi.Log)
		adminRouters.GET("/task/execute", api.TaskApi.Execute)
		adminRouters.GET("/task/stop", api.TaskApi.Stop)

		// 数据库操作
		adminRouters.GET("/mysql/table-list", api.MysqlApi.Tables)
		adminRouters.POST("/mysql/execute-sql", api.MysqlApi.Execute)

		// Redis缓存操作
		adminRouters.POST("/cache/redis-query", api.RedisApi.Query)

		// 店铺
		adminRouters.GET("/shop/list", api.ShopApi.List)
		adminRouters.POST("/shop/save", api.ShopApi.Save)
		adminRouters.GET("/shop/detail", api.ShopApi.Detail)

		// 终端
		adminRouters.POST("/terminal/change-status", api.TerminalApi.ChangeStatus)
		adminRouters.POST("/terminal/save", api.TerminalApi.Save)
		adminRouters.GET("/terminal/delete", api.TerminalApi.Delete)

		adminRouters.GET("/table/enable", api.TableApi.Enable)
		adminRouters.POST("table/save", api.TableApi.Save)
	}
}
