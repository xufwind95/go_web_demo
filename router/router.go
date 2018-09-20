package router

import (
	"net/http"

	_ "apiserver/docs"
	"go_web_demo/handler/sd"
	"go_web_demo/handler/user"
	"go_web_demo/router/middleware"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// 回复之前调用失败的接口的功能
	g.Use(gin.Recovery())
	// 使用路由器中自定义的中间件
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	// 使用外部传入的中间件
	g.Use(mw...)
	// 404 页面
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// 添加api文档
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// pprof router
	pprof.Register(g)

	// 登录
	g.POST("/login", user.Login)

	// 指定用户操作路由
	u := g.Group("/v1/user")
	u.Use(middleware.AuthMiddleware()) // 设置验证中间件
	{
		u.POST("", user.Create)
		u.DELETE("/:id", user.Delete)
		u.PUT("/:id", user.Update)
		u.GET("", user.List)
		u.GET("/:username", user.Get)
	}

	// 指定具体的路由,注意: 下面的大括号和上面的东西是隔开的，不是一个语句！
	svcd := g.Group("sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("ram", sd.RAMCheck)
	}

	return g
}
