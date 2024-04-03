package router

import (
	"github.com/gin-gonic/gin"

	"gocode/project/chatroom_api/user-web/api"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	//设置路由分组
	BaseRouter := Router.Group("base")
	{
		//设置URL 获取验证码, 将方法注册进路由
		BaseRouter.GET("captcha", api.GetCaptcha)
		//设置URL 获取手机短信, 将方法注册进路由
		BaseRouter.POST("send_sms", api.SendSms)
	}
}
