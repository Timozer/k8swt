package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Timozer/k8swt/common"
	"github.com/Timozer/k8swt/controllers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ReqId() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqId := strings.ToUpper(uuid.NewString())
		c.Set(common.HEADER_REQ_ID, reqId)
		c.Writer.Header().Set(common.HEADER_REQ_ID, reqId)
		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqId := c.GetString(common.HEADER_REQ_ID)
		logger := log.With().Str("ReqId", reqId).Logger()
		c.Set(common.CTX_LOGGER, &logger)
		c.Next()
	}
}

func InitLogger() {
	zerolog.TimestampFieldName = "T"
	zerolog.LevelFieldName = "L"
	zerolog.MessageFieldName = "M"

	logOut := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	logOut.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s: ", i)
	}

	// multi := zerolog.MultiLevelWriter(logOut, os.Stdout)
	multi := zerolog.MultiLevelWriter(logOut)

	log.Logger = zerolog.New(multi).With().Timestamp().Logger()
	// log.Logger = log.Output(logOut)
}

func main() {
	InitLogger()

	r := gin.Default()
	r.Use(ReqId(), Logger())
	r.Static("/js", "./dist/js")
	r.LoadHTMLGlob("./dist/*.html")

	r.GET("/", controllers.Index)
	r.GET("/ws", controllers.WsProcess)
	r.GET("/api/pods", controllers.Pods)

	r.Run(":8080")
}
