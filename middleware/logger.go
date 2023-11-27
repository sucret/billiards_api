package middleware

import (
	"billiards/pkg/log"
	"bytes"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"time"
)

// 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		bodyByte, _ := c.GetRawData()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyByte))

		c.Next()

		log.GetLogger().Info(
			c.Request.URL.Path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.Any("header", c.Request.Header),
			zap.ByteString("body", bodyByte),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", time.Since(start)),
		)
	}
}
