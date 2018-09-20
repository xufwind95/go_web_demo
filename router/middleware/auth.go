package middleware

import (
	"go_web_demo/handler"
	"go_web_demo/pkg/errno"
	"go_web_demo/pkg/token"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := token.ParseRequest(c); err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			// abort 这句不能少，不然会继续往下面走！
			c.Abort()
			return
		}

		c.Next()
	}
}
