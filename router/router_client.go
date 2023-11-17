package router

import (
	api "billiards/api/client"
	"github.com/gin-gonic/gin"
)

func setClientRoute(r *gin.Engine) {

	// code换session
	r.GET("/c/user/login", api.UserApi.Login)

	//clientRouters := r.Group("/c").Use(middleware.JWTAuth(service.AppGuardName), middleware.CheckAdminPermission(), gin.Logger(), middleware.CustomRecovery())
	//{
	//	clientRouters.GET("/user/1code2session", api.UserApi.Code2Session)
	//
	//}
}
