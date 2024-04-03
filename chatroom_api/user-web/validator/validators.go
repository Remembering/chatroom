package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// 自定义网页表单规则
// 使用正则表达式对手机号
func ValidateMobile(fl validator.FieldLevel) bool {
	//去url的多余字符
	mobile := fl.Field().String()
	//设置正则表达式
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if !ok {
		return ok
	}
	return ok
}
