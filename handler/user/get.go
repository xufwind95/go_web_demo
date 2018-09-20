package user

import (
	"github.com/gin-gonic/gin"
	. "go_web_demo/handler"
	"go_web_demo/model"
	"go_web_demo/pkg/errno"
)

func Get(c *gin.Context) {
	username := c.Param("username")

	user, err := model.GetUser(username)
	if err != nil {
		SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	SendResponse(c, nil, user)
}
