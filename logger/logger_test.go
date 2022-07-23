package logger

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"testing"
)

func TestLoggerToFile(t *testing.T) {
	router := gin.Default()
	router.Use(LoggerToFile())
	router.GET("/", func(context *gin.Context) {
		logger := getLogger()
		context.JSON(http.StatusOK, gin.H{
			"msg": "您好",
		})
		logger.Infof("测试测试")
	})
	router.Run(":8082")
}
