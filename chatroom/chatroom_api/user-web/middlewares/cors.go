package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
浏览器发现该请求是跨域请求,就会发送两个请求option请求和原本的请求
原本的请求服务器不会处理
*/
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		/*
			Access-Control-Allow-Origin：允许来自任何域的请求
			Access-Control-Allow-Headers：允许发送的HTTP头
			Access-Control-Allow-Methods：允许的HTTP方法
			Access-Control-Expose-Headers：允许从服务器返回的HTTP头
			Access-Control-Allow-Credentials：是否允许发送Cookie
		*/
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
