package user

import (
	. "go_web_demo/handler"
	"go_web_demo/model"
	"go_web_demo/pkg/errno"
	"go_web_demo/util"

	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/lexkong/log/lager"
)

// @Summary Add new user to the database
// @Description Add a new user
// @Tags user
// @Accept json
// @Produce json
// @Param user body user.CreateRequest true "Create a new user"
// @Success 200 {objects} user.CreateResponse "{"code":0,"message":"OK","data":{"username":"kong"}}"
// @Router /usr [post]
func Create(c *gin.Context) {
	log.Info("User Create function called.", lager.Data{"X-Request-Id": util.GetReqID(c)})
	var r CreateRequest

	// 获取 HTTP BODY 中的参数(json格式的可以这样搞)
	if err := c.Bind(&r); err != nil {
		SendResponse(c, errno.ErrBind, nil)
		return
	}

	u := model.UserModel{
		Username: r.Username,
		Password: r.Password,
	}

	if err := u.Validate(); err != nil {
		SendResponse(c, errno.ErrValidation, nil)
		return
	}

	if err := u.Encrypt(); err != nil {
		SendResponse(c, errno.ErrEncrypt, nil)
		return
	}

	if err := u.Create(); err != nil {
		SendResponse(c, errno.ErrDatabase, nil)
		return
	}

	rsp := CreateResponse{
		Username: r.Username,
	}

	SendResponse(c, nil, rsp)

	/*
		// 获取URL中的地址参数
		admin2 := c.Param("username")
		log.Infof("URL username: %s", admin2)

		//获取URL后面跟的键值对参数
		desc := c.Query("desc")
		log.Infof("URL key param desc: %s", desc)

		//获取HTTP HEADER中的参数
		contentType := c.GetHeader("Content-Type")
		log.Infof("Header Content-Type: %s", contentType)

		log.Debugf("username is: [%s], password is [%s]", r.Username, r.Password)
		if r.Username == "" {
			err := errno.New(
				errno.ErrUserNotFound,
				fmt.Errorf("username can not found in db: xx.xx.xx.xx"),
			).Add("This is add message.")
			log.Errorf(err, "Get an error")
			SendResponse(c, err, nil)
		}

		// 密码为空的，直接使用系统创建的错误
		if r.Password == "" {
			SendResponse(c, fmt.Errorf("password is empty"), nil)
		}

		rsp := CreateResponse{
			Username: r.Username,
		}
		SendResponse(c, nil, rsp)
	*/
}
