package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"gocode/project/chatroom_api/user-web/api"
	"gocode/project/chatroom_api/user-web/middlewares"
)

func InitUserRouter(Router *gin.RouterGroup) {
	//设置路由分组
	UserRouter := Router.Group("/user")

	//打印日志
	zap.S().Infof("配置用户相关URL")

	{ //设置URL以及中间件
		//输出所有用户,使用中间件并将方法注册进路由
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.IsAdminAuth(), api.GetUserList)
		//登录逻辑, 将方法注册进路由
		UserRouter.POST("pwd_login", api.PassWordLogin)
		//注册逻辑, 将方法注册进路由
		UserRouter.POST("register", api.Register)
	}
}
