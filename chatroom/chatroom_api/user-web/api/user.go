package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"gocode/project/chatroom_api/user-web/forms"
	"gocode/project/chatroom_api/user-web/global"
	"gocode/project/chatroom_api/user-web/global/response"
	"gocode/project/chatroom_api/user-web/middlewares"
	"gocode/project/chatroom_api/user-web/models"

	// "gocode/project/chatroom_api/user-web/initialize"
	"gocode/project/chatroom_api/user-web/proto"
)

// 将错误的多余字符去除
func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// 将grpc的code转换成http的状态码
func HandlerGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务信息不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})
			}
			return
		}
	}
}

func HanldeValidatorError(c *gin.Context, err error) {
	// 使用zh
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		//validator约束 字段不匹配的问题
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
}

// 获取所有用户列表的grpc客户端
func GetUserList(ctx *gin.Context) {
	//创建于grpc服务器的连接
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d",
		global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 [用户服务] 失败", "msg", err.Error())
	}
	defer userConn.Close()

	//ctx的上下文 设置一个实例 获取“caims” 的值
	claims, _ := ctx.Get("claims")
	//值底层类型是interfacep{} 所以要断言
	currenUser := claims.(*models.CustomClaims)

	//日志输出
	zap.S().Infof("访问用户: %d", currenUser.ID)

	//生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)

	//获取url的query值, 若没有则使用默认值
	pn := ctx.DefaultQuery("pn", "0")
	// 字符串转整型
	pnInt, _ := strconv.Atoi(pn)
	pSize := ctx.DefaultQuery("psize", "10")
	pSizeInt, _ := strconv.Atoi(pSize)

	//使用grpc服务器里的GetUserList函数,并接收返回数据rsp
	rsp, err := userSrvClient.GetUserList(global.Ctx, &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 [用户列表] 失败")
		HandlerGrpcErrorToHttp(err, ctx)
		return
	}
	//实例一个接收任意类型的数组
	result := make([]interface{}, 0)
	// 将grpc服务器的rsp里面的Date数据进行保存
	for _, value := range rsp.Date {
		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			// Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)).Format("2006-01-02"),
			BirthDay: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}

		result = append(result, user)
	}
	//网页输出
	ctx.JSON(http.StatusOK, result)
}

// 密码登陆的grpc客户端
func PassWordLogin(c *gin.Context) {
	//1. 表单验证
	passWordLoginForm := forms.PassWordLoginForm{}
	/*
		1.匹配表单的类型是什么(用什么方式进行传递的)
		2.确认类型确认tag是否符合约束
		3.符合就把request.body的相应内容加载进struc里
	*/
	if err := c.ShouldBind(&passWordLoginForm); err != nil {
		HanldeValidatorError(c, err)
		return
	}

	//验证码 ture 验证码用一次就过期
	if !store.Verify(passWordLoginForm.CaptchaId, passWordLoginForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	//拨号连接用户grpc服务
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d",
		global.ServerConfig.UserSrvInfo.Host,  //ip
		global.ServerConfig.UserSrvInfo.Port), //端口
		grpc.WithTransportCredentials(insecure.NewCredentials())) //不安全验证
	if err != nil {
		zap.S().Errorw("[PassWordLogin] 连接 [用户服务] 失败", "msg", err.Error())
		return
	}
	defer userConn.Close()
	//生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)

	//登陆的逻辑

	//1. 查询用户是否存在
	userRsp, err := userSrvClient.GetUserByMobile(global.Ctx, &proto.MobileRequest{
		Mobile: passWordLoginForm.Mobile,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusBadRequest, gin.H{
					"mobile": "用户不存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, map[string]string{
					"mobile": "登录失败",
				})
			}
			return
		}
	}
	//2. 密码是否正确
	passRsp, _ := userSrvClient.CheckPassword(global.Ctx, &proto.PassWordCheckInfo{
		Password:          passWordLoginForm.PassWord,
		EncryptedPassword: userRsp.PassWord,
	})
	if !passRsp.Success {
		c.JSON(http.StatusBadRequest, map[string]string{
			"password": "密码错误",
		})
		return
	}
	//密码正确逻辑
	//1. 生成token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(userRsp.Id),
		NickName:    userRsp.NickName,
		AuthorityId: uint(userRsp.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //签名的生效时问
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //签名失效时间
			Issuer:    "LMXLJLX",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}
	//网页输出
	c.JSON(http.StatusOK, gin.H{
		"Id":          userRsp.Id,
		"nick_name":   userRsp.NickName,
		"role":        userRsp.Role,
		"token":       token,
		"expire_time": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}

// 用户注册的grpc客户端
func Register(c *gin.Context) {
	//用户注册表单类型实例
	registerForm := forms.RegisterForm{}
	//验证是否符合bind的要求
	if err := c.ShouldBind(&registerForm); err != nil {
		HanldeValidatorError(c, err)
		return
	}

	// 校验验证码
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})
	//向redis里查询该手机的验证码是否匹配
	value, err := rdb.Get(global.Ctx, registerForm.Mobile).Result()
	if err == redis.Nil || value != registerForm.Code {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	}

	//创建于grpc服务器的连接
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d",
		global.ServerConfig.UserSrvInfo.Host,
		global.ServerConfig.UserSrvInfo.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[Register] 连接 [用户服务] 失败", "msg", err.Error())
	}
	defer userConn.Close()
	userSrvClient := proto.NewUserClient(userConn)
	//创建用户
	user, err := userSrvClient.CreateUser(global.Ctx, &proto.CreatUserInfo{
		NickName: registerForm.Mobile,
		PassWord: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorf("[Register] 查询 [新建用户] 失败:%s", err.Error())
		HandlerGrpcErrorToHttp(err, c)
		return
	}

	//注册自动登录

	//密码正确逻辑
	//1. 生成token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //签名的生效时问
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //签名失效时间
			Issuer:    "LMXLJLX",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}
	//网页输出
	c.JSON(http.StatusOK, gin.H{
		"Id":          user.Id,
		"nick_name":   user.NickName,
		"token":       token,
		"expire_time": (time.Now().Unix() + 60*60*24*30) * 1000,
	})
}
