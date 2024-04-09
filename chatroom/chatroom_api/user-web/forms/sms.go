// 表单
package forms

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` //手机号码格式规范,自定义validator
	Type   uint   `form:"type" json:"type" binding:"required,oneof=1 2"`  //1:register 2:login
	//1. 注册发送短信验证码和动态验证码登录发送验证码
}
