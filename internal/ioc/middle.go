package ioc

import (
	"classroom-analysis/internal/web/middleware"

	"github.com/gin-gonic/gin"
)

// 创建全局JWT中间件实例
var jwtMiddleware gin.HandlerFunc

func init() {
	// 初始化JWT中间件，忽略登录路径
	jwtMiddleware = middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/teacher/login").
		IgnorePaths("/teacher/register").
		Build()
}

func GetJWTMiddleware() gin.HandlerFunc {
	return jwtMiddleware
}
func InitMiddlewares(loggers []gin.HandlerFunc) []gin.HandlerFunc {
	// 组合现有的 loggers 和 你的 jwtMiddleware
	res := []gin.HandlerFunc{
		gin.Recovery(),
		gin.Logger(),
	}

	// 把传入的 loggers（如果有的话）加入进来
	res = append(res, loggers...)

	// 把 JWT 中间件加入进来
	if jwtMiddleware != nil {
		res = append(res, jwtMiddleware)
	}

	return res
}
