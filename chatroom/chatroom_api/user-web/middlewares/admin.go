package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gocode/project/chatroom_api/user-web/models"
)

// 是否是管理员的中间件
func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, _ := ctx.Get("claims")
		currenUser := claims.(*models.CustomClaims)
		if currenUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
