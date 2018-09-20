package user

import (
	. "go_web_demo/handler"
	"go_web_demo/model"
	"go_web_demo/pkg/auth"
	"go_web_demo/pkg/errno"
	"go_web_demo/pkg/token"

	"github.com/gin-gonic/gin"
)

// 用户登录功能
func Login(c *gin.Context) {
	// 将传入写入用户
	var u model.UserModel
	if err := c.Bind(&u); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	// 判定用户是否存在
	d, err := model.GetUser(u.Username)
	if err != nil {
		SendResponse(c, errno.ErrUserNotFound, nil)
		return
	}

	// 对比密码是否匹配
	if err := auth.Compare(d.Password, u.Password); err != nil {
		SendResponse(c, errno.ErrPasswordIncorrect, nil)
		return
	}

	// 生成token
	t, err := token.Sign(c, token.Context{ID: d.Id, Username: d.Username}, "")
	if err != nil {
		SendResponse(c, errno.ErrToken, nil)
		return
	}

	// 发送token
	SendResponse(c, nil, model.Token{Token: t})
}
