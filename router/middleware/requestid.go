package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

// 所有请求设置 Request Id
func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取 X-Requested-Id 便于后续验证
		requestId := c.Request.Header.Get("X-Request-Id")

		if requestId == "" {
			u4, _ := uuid.NewV4()
			requestId = u4.String()
		}

		c.Set("X-Request-Id", requestId)

		c.Writer.Header().Set("X-Request-Id", requestId)
		c.Next()
	}
}
