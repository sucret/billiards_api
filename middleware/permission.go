package middleware

import (
	"billiards/response"
	"billiards/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 检测访问权限
func CheckAdminPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		userId, _ := strconv.Atoi(c.GetString("userId"))

		if !service.AdminService.CheckAdminPermission(userId, path) {
			response.PermissionDenied(c)
			c.Abort()
			return
		}
	}
}
