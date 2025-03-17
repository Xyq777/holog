package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/ncuhome/holog"
)

func main() {
	logger := holog.NewLogger("test-service", holog.WithFileWriter(&lumberjack.Logger{
		Filename:   "./zap.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}))
	r := gin.New()
	r.Use(holog.HologGinRequestLogging(logger))

	r.GET("/ping", func(c *gin.Context) {
		fmt.Println("Received /ping request")
		time.Sleep(500 * time.Millisecond)

		err := errors.New("test error")
		c.Error(err)

		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	})

	r.Run(":8080")
}
