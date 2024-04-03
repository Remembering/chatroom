package initialize

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"

	"gocode/project/chatroom_api/user-web/global"
)

// 初始化翻译器并设置中文的翻译器
func InitTrans(locale string) (err error) {
	//修改gin框架中的validatio引擎熟悉,实现定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//注册一个获取json的tag的自定义方法
		v.RegisterTagNameFunc(func(field reflect.StructField) string {
			//分割成最多n个字符串
			name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
			// json tag 中 "-" 代表不处理
			if name == "-" {
				return name
			}
			return name

		})
		zhT := zh.New() // 中文翻译器
		enT := en.New() // 英文翻译器
		//第一个参数时备用语言环境, 后面参数是应该支持的语言环境
		uni := ut.New(enT, zhT, zhT)
		var ok bool
		global.Trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTransloator(%s)", locale)
		}
		switch locale {
		case "en":
			en_translations.RegisterDefaultTranslations(v, global.Trans)
		case "zh":
			zh_translations.RegisterDefaultTranslations(v, global.Trans)
		default:
			en_translations.RegisterDefaultTranslations(v, global.Trans)
		}
		return
	}
	return
}
