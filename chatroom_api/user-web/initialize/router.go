package initialize

import (
	"github.com/gin-gonic/gin"

	"gocode/project/chatroom_api/user-web/middlewares"
	"gocode/project/chatroom_api/user-web/router"
)

func Routers() *gin.Engine {
	//创建路由实例, 使用默认的路由
	Router := gin.Default()

	//配置跨域
	Router.Use(middlewares.Cors())

	//设置路由分组
	ApiGroup := Router.Group("/u/v1")

	//以下路由都是该分组下的
	//初始化用户路由
	router.InitUserRouter(ApiGroup)
	//初始化基本路由
	router.InitBaseRouter(ApiGroup)
	return Router
}
