package main

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"gocode/project/chatroom_api/user-web/global"
	"gocode/project/chatroom_api/user-web/initialize"
	myvalidator "gocode/project/chatroom_api/user-web/validator"
)

func main() {
	//1. 初始化日志
	initialize.InitLogger()
	//2. 初始化配置文件
	initialize.InitConfig()
	//3. 初始化routers
	Router := initialize.Routers()
	//4. 初始化翻译validators
	if err := initialize.InitTrans("zh"); err != nil {
		fmt.Println("初始化翻译错误")
		zap.S().Errorf("[initialize.InitTrans] 初始化 错误", "msg", err.Error())
		return
	}
	//5. 注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//注册validate实例
		v.RegisterValidation("mobile", myvalidator.ValidateMobile)
		//注册实例的错误返回
		v.RegisterTranslation("mobile", global.Trans,
			func(ut ut.Translator) error {
				return ut.Add("mobile", "{0} 非法的手机号码!", true)
			},
			func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T("mobile", fe.Field())
				return t
			})
	}
	/*
		1. S()可以获取一个全局的sugar，可以让我们自己设置一个全局的Logger
		2. 日志是分级别的，debug, info, warn, error, fetal 低到高
		当设置高级别,低级别的日志不会输出
		3. S函数和L函数很有用，提供了一个全局的安全访问1ogger的途径
	*/
	zap.S().Debugf("启动服务, 端口:%d", global.ServerConfig.Port)
	//网页服务器监听
	if err := Router.Run(fmt.Sprintf(":%d", global.ServerConfig.Port)); err != nil {
		zap.S().Panic("启动服务失败:", err.Error())
	}
}
