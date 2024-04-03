package api

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"gocode/project/chatroom_api/user-web/forms"
	"gocode/project/chatroom_api/user-web/global"
)

// 发送短信
func SendSms(ctx *gin.Context) {

	//1. 表单验证
	sendSmsForm := forms.SendSmsForm{}
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		HanldeValidatorError(ctx, err)
		return
	}

	client, err := CreateClient(tea.String(global.ServerConfig.AliSmsInfo.ApiKey), tea.String(global.ServerConfig.AliSmsInfo.ApiSecret))
	if err != nil {
		zap.S().Errorw("[SMS] [SendSms] error", "msg", err.Error())
	}

	rand.Seed(time.Now().UnixNano()) //每次重置初始值
	smsCode := fmt.Sprintf("%04d", rand.Intn(10000))
	smsCodeFormat := fmt.Sprintf("{\"code\":%s}", smsCode)

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		SignName:     tea.String("网络聊天室"),
		TemplateCode: tea.String("SMS_465375633"),
		PhoneNumbers: tea.String(sendSmsForm.Mobile),
	}
	//设置动态随机验证码
	sendSmsRequest.TemplateParam = tea.String(smsCodeFormat)

	runtime := &util.RuntimeOptions{}
	{
		runtime.SetReadTimeout(10000)   //  读取超时
		runtime.SetConnectTimeout(5000) // 连接超时
		runtime.SetAutoretry(true)      // 是否自动重试
	}

	// 获取响应对象
	_, err = client.SendSmsWithOptions(sendSmsRequest, runtime)
	if err != nil {
		zap.S().Errorw("[SMS] [SendSms] error", "msg", err.Error())
	}

	//连接ridis 保存验证码
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	//设置手机号与验证码的键值对
	rdb.Set(global.Ctx, sendSmsForm.Mobile, smsCode, time.Duration(global.ServerConfig.RedisInfo.Expire)*time.Second)
	//向网页输出
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "发送成功",
	})
}

func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		Endpoint:        tea.String("dysmsapi.aliyuncs.com"),
	}
	// Endpoint 请参考 https://api.aliyun.com/product/Dysmsapi
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

// 按with生成多少为的验证码
func GenerateSmsCode(witdh int) string {
	var n int = 1
	for i := 0; i < witdh; i++ {
		n *= 10
	}
	format := "%0" + fmt.Sprintf("%d", witdh) + "d"
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf(format, rand.Intn(n))
}
