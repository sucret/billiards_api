package middleware

import (
	"github.com/gin-gonic/gin"
)

// 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		//start := time.Now()
		//
		//bodyByte, _ := c.GetRawData()
		//c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyByte))

		c.Next()

		//log.GetLogger().Info(
		//	c.Request.URL.Path,
		//	zap.Int("status", c.Writer.Status()),
		//	zap.String("method", c.Request.Method),
		//	zap.String("path", c.Request.URL.Path),
		//	zap.String("query", c.Request.URL.RawQuery),
		//	zap.ByteString("body", bodyByte),
		//	zap.String("ip", c.ClientIP()),
		//	zap.String("user-agent", c.Request.UserAgent()),
		//	zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
		//	zap.Duration("cost", time.Since(start)),
		//)
	}
}
