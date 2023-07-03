package main

import (
	"flag"
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

func TermVars(args *Arguments) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(common.TERM_SHELLS, *args.Shells)
		c.Next()
	}
}

func InitLogger(args *Arguments) {
	zerolog.TimestampFieldName = "T"
	zerolog.LevelFieldName = "L"
	zerolog.MessageFieldName = "M"

	logOut := zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}
	logOut.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s: ", i)
	}

	// multi := zerolog.MultiLevelWriter(logOut, os.Stdout)
	multi := zerolog.MultiLevelWriter(logOut)

	log.Logger = zerolog.New(multi).With().Timestamp().Logger().
		Level(zerolog.Level(*args.LogLevel))
	// log.Logger = log.Output(logOut)
}

type Arguments struct {
	Dev      *bool
	Port     *string
	LogLevel *int
	Shells   *string
}

var (
	gArgs = &Arguments{}
)

func init() {
	gArgs.Dev = flag.Bool("dev", false, "dev mode")
	gArgs.Port = flag.String("port", "8080", "http port")
	gArgs.LogLevel = flag.Int("loglevel", 1, "log level: -1 trace, 0 debug, 1 info, 2 warn, 3 error, 4 fatal, 5 panic, 6 nolevel, 7 disabled")
	gArgs.Shells = flag.String("shells", "/bin/bash:/bin/sh", "the shells to use, seprated by colon")
}

func main() {
	flag.Parse()

	InitLogger(gArgs)

	if !*gArgs.Dev {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(ReqId(), Logger())

	if *gArgs.Dev {
		r.Static("/js", "web/dist/js")
		r.LoadHTMLGlob("web/dist/*.html")
	} else {
		r.Static("/js", "./dist/js")
		r.LoadHTMLGlob("./dist/*.html")
	}

	r.GET("/", controllers.Index)
	r.GET("/ws", TermVars(gArgs), controllers.WsProcess)
	r.GET("/api/pods", controllers.Pods)

	r.Run(":" + *gArgs.Port)
}
