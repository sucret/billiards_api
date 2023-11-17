package middleware

import (
	"billiards/pkg/config"
	"billiards/response"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

func CustomRecovery() gin.HandlerFunc {
	logConfig := config.GetConfig().Log
	return gin.RecoveryWithWriter(
		&lumberjack.Logger{
			Filename:   logConfig.RootDir + "/" + logConfig.ServerErrorLog,
			MaxSize:    logConfig.MaxSize,
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge,
			Compress:   logConfig.Compress,
		},
		response.ServerError)
}
