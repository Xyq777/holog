package holog

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

func HologGinRequestLogging(logger *logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		startTime := time.Now()

		reqBodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBodyBytes))
		reqBodyString := string(reqBodyBytes)

		c.Next()

		latency := fmt.Sprintf("%.2fms", float32(time.Since(startTime).Microseconds())/1000.0)
		if c.Writer.Status() >= 400 {
			logger.Error("Request",
				"status", c.Writer.Status(),
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"latency", latency,
				"req_headers", c.Request.Header,
				"ip", c.ClientIP(),
				"req_body", reqBodyString,
				"error", c.Errors.String(),
			)
		} else {
			logger.Info("Request",
				"status", c.Writer.Status(),
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"latency", latency,
				"req_headers", c.Request.Header,
				"ip", c.ClientIP(),
				"req_body", reqBodyString,
			)
		}
	}
}
