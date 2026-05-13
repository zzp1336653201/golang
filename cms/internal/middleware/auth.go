package middleware

import (
	"net/http"
	"strings"

	"cms/internal/model"
	"cms/internal/service"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.Error("请先登录"))
			return
		}
		
		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.Error("token格式错误"))
			return
		}
		
		token := parts[1]
		
		// 验证token
		user, err := userService.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, model.Error(err.Error()))
			return
		}
		
		// 将用户信息存入Context
		c.Set("user_id", user.ID)
		c.Set("user", user)
		
		c.Next()
	}
}

// GetUserID 从Context获取用户ID
func GetUserID(c *gin.Context) uint {
	if id, exists := c.Get("user_id"); exists {
		return id.(uint)
	}
	return 0
}
