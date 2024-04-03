package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

// 实例一个base64编码的store
var store = base64Captcha.DefaultMemStore

// 获取验证码
func GetCaptcha(ctx *gin.Context) {
	// driver := base64Captcha.NewDriverDigit(64, 240, 5, 0.7, 80)
	//使用默认的设置 长 宽 高 扭曲因子 数字验证码
	driver := base64Captcha.DefaultDriverDigit
	// 实例一个验证码对象
	c := base64Captcha.NewCaptcha(driver, store)
	//生成一个验证码 返回 id b64s编码的图像 数字验证码
	id, b64s, answer, err := c.Generate()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	//响应返回
	ctx.JSON(http.StatusOK, gin.H{
		"id":     id,
		"b64s":   b64s,
		"answer": answer,
	})
}
