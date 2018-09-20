package user

import (
	"github.com/gin-gonic/gin"
	. "go_web_demo/handler"
	"go_web_demo/model"
	"go_web_demo/pkg/errno"
	"strconv"
)

func Delete(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	if err := model.DeleteUser(uint64(userId)); err != nil {
		SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	SendResponse(c, nil, nil)
}
