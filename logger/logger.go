package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
)

// 检查并创建日志文件夹
func createFolder(logFilePath string) {
	if err := os.MkdirAll(logFilePath, 0777); err != nil {
		fmt.Println("文件夹创建失败")
		panic(err)
	}
}

// 创建 打开日志文件
func createLogFile(logFilePath, filenameFormat string) *os.File {
	logFileName := filenameFormat + ".log"
	filename := path.Join(logFilePath, logFileName)

	// 检查是否能够成功创建日志文件
	checkFile := func(filename string) {
		if _, err := os.Stat(filename); err != nil {
			if _, err := os.Create(filename); err != nil {
				fmt.Println("打开文件失败")
				panic(err)
			}
		}
	}

	checkFile(filename)
	src, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	return src
}

// 获取日志写入的logger
func getLogger() *logrus.Logger {
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}

	// 创建日志文件夹
	createFolder(logFilePath)
	filenameFormat := time.Now().Format("2006-01-02")
	src := createLogFile(logFilePath, filenameFormat)

	logger := logrus.New()

	// 输出到控制台
	writers := []io.Writer{
		src,
		os.Stdout,
	}
	// 同时写文件和屏幕
	fileAndStdoutWriter := io.MultiWriter(writers...)
	logger.SetOutput(fileAndStdoutWriter)

	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return logger
}

func LoggerToFile() gin.HandlerFunc {
	logger := getLogger()
	fmt.Println("使用middle ware")
	return func(context *gin.Context) {
		fmt.Println("收到请求")

		startTime := time.Now()

		context.Next()
		// 处理请求
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)

		reqMethod := context.Request.Method
		// 请求路由
		reqUri := context.Request.RequestURI

		// 状态码
		statusCode := context.Writer.Status()

		// 日志格式
		logger.Infof("| %3d | %13v | %15s | %s",
			statusCode,
			latencyTime,
			reqMethod,
			reqUri,
		)
	}
}
